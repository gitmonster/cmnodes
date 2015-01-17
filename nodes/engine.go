package nodes

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/librato"
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
	TypeMap[GetNodeTypeName(node)] = fn
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
func (e *Engine) NodeExists(scope, nodeId string) (bool, error) {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	if cnt, err := coll.FindId(nodeId).Count(); err != nil {
		return false, err
	} else if cnt == 1 {
		return true, nil
	}

	return false, nil
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

	return false, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) MoveNode(scope, srcNodeType, srcNodeId, targetNodeId string) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	// check if source node exists
	if ex, err := e.NodeExists(scope, srcNodeId); err != nil {
		return err
	} else if !ex {
		return fmt.Errorf("MoveNode::Source node %s doesn't exist.", targetNodeId)
	}

	// check if target node exists
	if ex, err := e.NodeExists(scope, targetNodeId); err != nil {
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
func (e *Engine) CreateInstanceByType(nodeType string) (Node, error) {
	if fn, ok := TypeMap[nodeType]; !ok {
		return nil, fmt.Errorf("CreateInstanceByType::Type %s is not registered yet.", nodeType)
	} else {

		session, coll := e.GetMgoSession(PROTOS_SCOPE)
		defer session.Close()

		query := coll.FindId(nodeType)
		if cnt, err := query.Count(); err != nil {
			return nil, err
		} else if cnt == 0 {
			return nil, fmt.Errorf("CreateInstanceByType::Prototype for NodeType %s is not available", nodeType)
		}

		node := fn(e)
		if err := query.One(node); err != nil {
			return nil, err
		}

		node.NewObjectId()
		return node, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNewNode(scope, parentId, name, nodeType string) (Node, error) {
	if node, err := e.CreateInstanceByType(nodeType); err != nil {
		return nil, err
	} else {
		node.SetParentId(parentId)
		node.SetName(name)

		session, coll := e.GetMgoSession(scope)
		defer session.Close()

		if err := coll.Insert(node); err != nil {
			return nil, err
		}

		return node, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) RegisterAllRoutes(scope string, router *mux.Router) error {
	session, coll := e.GetMgoSession(scope)
	defer session.Close()

	iter := coll.Find(bson.M{"rr": true}).Iter()
	if err := iter.Err(); err != nil {
		return err
	}

	var node interface{}
	for iter.Next(&node) {
		if n, ok := node.(Node); !ok {
			return fmt.Errorf("RegisterAllRoutes::Could not assert child to Node type")
		} else {
			n.RegisterRoute(router)
		}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) AssembleRouteFor(scope, nodeId string) string {
	return EMPTY_STRING
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
	if err := e.RegisterAllRoutes(SYSTEM_SCOPE, sysRouter); err != nil {
		return err
	}

	e.negroni = negroni.Classic()
	e.negroni.UseHandler(mainRouter)
	e.negroni.Run(connection)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadProtoContent(typeName, section string) (string, error) {

	name := fmt.Sprintf("%s.%s", typeName, section)
	pt := path.Join(e.StartupDir, "scopes", "protos", name)

	if _, err := os.Stat(pt); err != nil {
		if os.IsNotExist(err) {
			return EMPTY_STRING, nil
		} else {
			return EMPTY_STRING, err
		}
	}

	if buf, err := ioutil.ReadFile(pt); err != nil {
		return EMPTY_STRING, err
	} else {
		return string(buf), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) ImportPrototypes(force bool) error {
	e.Logger.Infof("Import prototypes")

	session, coll := e.GetMgoSession(PROTOS_SCOPE)
	defer session.Close()

	if force {
		coll.RemoveAll(nil)
	}

	for tp, fn := range TypeMap {
		query := coll.FindId(tp)
		if cnt, err := query.Count(); err != nil {
			return err
		} else if cnt == 0 {
			e.Logger.Infof("Import prototype for %s", tp)

			node := fn(e)
			node.SetObjectId(tp)

			if cont, err := e.LoadProtoContent(tp, "edit"); err != nil {
				return err
			} else {
				node.SetEditTemplate(cont)
			}

			if err := coll.Insert(node); err != nil {
				return err
			}
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) EnsurePrototypes() error {

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) EnsureSystem() error {
	if err := e.EnsurePrototypes(); err != nil {
		return err
	}

	_, err := e.CreateNewNode(SYSTEM_SCOPE, EMPTY_STRING,
		"System", NODETYPE_SITE)

	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	return fn(e)
}

////////////////////////////////////////////////////////////////////////////////
func NewEngine(config *NodesConfig) (*Engine, error) {
	eng := Engine{Config: config}
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

	eng.EnsureSystem()
	return &eng, nil
}
