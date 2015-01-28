package nodes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/librato"
	"gopkg.in/validator.v2"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

////////////////////////////////////////////////////////////////////////////////
type EngineFunc func(engine *Engine) error

type Engine struct {
	Loggable
	MongoSessionProvider
	negroni    *negroni.Negroni
	Config     *NodesConfig
	SystemCriteria *CriteriaSet
	StartupDir string
}

var (
	TypeMap map[string]NewNodeFunc
)

////////////////////////////////////////////////////////////////////////////////
func init() {
	TypeMap = make(map[string]NewNodeFunc)
}

type EnumFunc func(node Node) error
type NewNodeFunc func(engine *Engine) Node

////////////////////////////////////////////////////////////////////////////////
func RegisterNodeType(node interface{}, fn NewNodeFunc) {
	TypeMap[GetTypeName(node)] = fn
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) EnumerateChilds(scope, nodeId string, fn EnumFunc) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	iter := coll.Find(bson.M{"p": nodeId}).Iter()
	if err := iter.Err(); err != nil {
		return err
	}

	var node interface{}
	for iter.Next(&node) {
		if n, ok := node.(Node); !ok {
			return fmt.Errorf("EnumerateChilds::Could not assert child to Node type")
		} else {
			if err := fn(n); err != nil {
				return err
			}
		}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CheckNodeIsNoChild(scope, nodeId string, parentNodeId string) (bool, error) {
	if nodeId == parentNodeId {
		return false, nil
	}

	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	iter := coll.Find(bson.M{"p": nodeId}).Iter()
	if err := iter.Err(); err != nil {
		return false, err
	}

	node := NodeBase{}
	for iter.Next(&node) {
		if nodeId == node.Id {
			return false, nil
		}
	}

	if err := iter.Close(); err != nil {
		return false, err
	}

	return true, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CheckIsChildAllowed(scope, parentId, nodeType string) (bool, error) {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	query := coll.FindId(parentId)

	var node interface{}
	if err := query.One(&node); err != nil {
		return false, err
	} else {
		if n, ok := node.(Node); !ok {
			return false, fmt.Errorf("CheckIsChildAllowed::Could not assert child to Node type")
		} else {
			return n.IsChildAllowed(nodeType), nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) MoveNode(scope, srcNodeType, srcNodeId, targetNodeId string) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	// check if source node exists

	crit := NewCriteria(scope).WithId(srcNodeId)
	if ex, err := e.NodeExists(crit); err != nil {
		return err
	} else if !ex {
		return fmt.Errorf("MoveNode::Source node %s doesn't exist.", targetNodeId)
	}

	crit.WithObjectId(targetNodeId)
	// check if target node exists
	if ex, err := e.NodeExists(crit); err != nil {
		return err
	} else if !ex {
		return fmt.Errorf("MoveNode::Target node %s doesn't exist.", targetNodeId)
	}

	if ok, err := e.CheckNodeIsNoChild(scope, srcNodeId, targetNodeId); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("MoveNode::Source node %s is already a child of Target node %s.", srcNodeId, targetNodeId)
	}

	if ok, err := e.CheckIsChildAllowed(scope, targetNodeId, srcNodeType); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("MoveNode::Source node %s is not allowed as child of Target node %s.", srcNodeId, targetNodeId)
	}

	set := bson.M{"$set": bson.M{
		"p": targetNodeId, //set parent to targetNodeId
		"o": 0,            //set order to 0
	}}

	if err := coll.UpdateId(srcNodeId, set); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) RemoveNode(scope, nodeId string) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	// delete childs
	iter := coll.Find(bson.M{"p": nodeId}).Iter()
	if err := iter.Err(); err != nil {
		return err
	}

	node := NodeBase{}
	for iter.Next(&node) {
		e.RemoveNode(scope, node.Id)
	}

	if err := iter.Close(); err != nil {
		return err
	}

	// delete node
	if err := coll.RemoveId(nodeId); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateInstanceByType(nodeType string, abortNoPrototype bool) (Node, error) {
	if fn, ok := TypeMap[nodeType]; !ok {
		return nil, fmt.Errorf("CreateInstanceByType::Type %s is not registered yet.", nodeType)
	} else {
		session, coll := e.GetMgoSession(SYSTEM_SCOPE)
		defer session.Close()

		crit := NewCriteria(SYSTEM_SCOPE).
			WithParentId(OBJECTID_SYSTEM_PROTOTYPES).
			WithProtoNodeType(nodeType)

		if ex, err := e.NodeExists(crit); err != nil {
			return nil, err
		} else if abortNoPrototype && !ex {
			return nil, fmt.Errorf("CreateInstanceByType::Prototype for NodeType %s is not available, aborting", nodeType)
		} else {
			node := fn(e)
			if ex {
				if err := coll.Find(crit.GetSelector()).One(node); err != nil {
					return nil, err
				}
			}

			node.NewObjectId()
			return node, nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNode(crit *Criteria, abortNoPrototype bool) (Node, error) {
	if err := validator.Validate(crit); err != nil {
		return nil, err
	}

	if node, err := e.CreateInstanceByType(crit.GetNodeType(), abortNoPrototype); err != nil {
		return nil, err
	} else {
		session, coll := e.GetMgoSession(crit.GetScope())
		defer session.Close()

		if err := node.Apply(crit); err != nil {
			return nil, err
		}
		if _, err := coll.UpsertId(node.GetObjectId(), node); err != nil {
			return nil, err
		}

		return node, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) RegisterRoutesForScope(scope string, router *mux.Router) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	iter := coll.Find(bson.M{"rr": true}).Iter()
	if err := iter.Err(); err != nil {
		return err
	}

	var node interface{}
	for iter.Next(&node) {
		if n, ok := node.(Node); !ok {
			return fmt.Errorf("RegisterRoutesForScope::Could not assert child to Node type")
		} else{
			if r, err := e.AssembleRouteFor(scope, n.GetObjectId()); err != nil {
				return err
			} else {
				n.RegisterRoute(r, router)
			}
		}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) AssembleRouteFor(scope, nodeId string) (string, error) {
	return EMPTY_STRING, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Serve() error {
	if cc, err := e.Config.GetConnectionConfig(); err != nil {
		return err
	} else {
		connection := cc["connection"].(string)
		if err := e.Startup(connection); err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Startup(connection string) error {

	mainRouter := mux.NewRouter().StrictSlash(false)
	systemRouter := mux.NewRouter()

	mainRouter.PathPrefix("/nodes").Handler(negroni.New(
		negroni.NewRecovery(),
		//negroni.HandlerFunc(MyMiddleware),
		negroni.NewLogger(),
		negroni.Wrap(systemRouter),
	))

	sysRouter := systemRouter.PathPrefix("/nodes").Subrouter()
	if err := e.RegisterRoutesForScope(SYSTEM_SCOPE, sysRouter); err != nil {
		return err
	}

	e.negroni = negroni.Classic()
	e.negroni.UseHandler(mainRouter)
	e.negroni.Run(connection)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) EnsureNodeExists(crit *Criteria, abortNoPrototype bool) error {
	if ex, err := e.NodeExists(crit); err != nil {
		return err
	} else if !ex {
		if _, err := e.CreateNode(crit, abortNoPrototype); err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CheckSystemIntegrity() error {
	e.Logger.Infof("Check system integrity...")

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) ImportSystem(force bool, filePath string) error {
    if filePath = EMPTY_STRING{
        filePath = path.Join(e.StartupDir, "system.toml")
    }

    set := e.SystemCriteria
	if err := set.LoadFromFile(filePath); err != nil {
		return err
	}

	if err := set.Ensure(force,e); err != nil {
		return err
	}
    return nil
}


////////////////////////////////////////////////////////////////////////////////
func (e *Engine) NodeExists(crit *Criteria) (bool, error) {
	session, coll := e.GetMgoSession(crit.GetScope())
	defer session.Close()

	query := coll.FindId(crit.GetSelector())
	if cnt, err := query.Count(); err != nil {
		return false, err
	} else if cnt > 0 {
		return true, nil
	}

	return false, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	return fn(e)
}

////////////////////////////////////////////////////////////////////////////////
func NewEngine(config *NodesConfig) (*Engine, error) {
	eng := Engine{Config: config}
    eng.SystemCriteria = make(CriteriaSet)
	eng.SetLogger(eng.NewLogger("engine"))

	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return nil, err
	} else {
		eng.StartupDir = dir
	}

	if lConfig, err := config.GetLibratoConfig(); err != nil {
		return nil, err
	} else {

		go librato.Librato(metrics.DefaultRegistry,
			lConfig["duration"].(time.Duration),
			lConfig["email"].(string),
			lConfig["apitoken"].(string),
			lConfig["hostname"].(string),
			[]float64{0.95},  // precentiles to send
			time.Millisecond, // time unit
		)
	}

	if mConf, err := config.GetMongoConfig(); err != nil {
		return nil, err
	} else {

		if session, err := mgo.Dial(mConf["host"].(string)); err != nil {
			return nil, fmt.Errorf("NewEngine::Init error:: Mongo Session could not be initialized. :: %s", err.Error())
		} else {
			session.SetMode(mgo.Monotonic, true)
			eng.Session = session
		}
	}

	if err := eng.CheckSystemIntegrity(); err != nil {
		return nil, err
	}
	return &eng, nil
}
