package nodes

import (
	"fmt"
	"net/http"
	"reflect"
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
	negroni *negroni.Negroni
	Config  *NodesConfig
}

var (
	TypeMap map[string]reflect.Type
)

////////////////////////////////////////////////////////////////////////////////
func init() {
	TypeMap = make(map[string]reflect.Type)
}

////////////////////////////////////////////////////////////////////////////////
func RegisterNodeType(node interface{}) {
	t := reflect.TypeOf(node)
	TypeMap[t.Name()] = t
}

type EnumFunc func(node Node) error

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

    node := interface{}
	if , err := query.One(&node); err != nil {
		return false, err
	} else
	    if n, ok := node.(Node); !ok {
			return fmt.Errorf("CheckIsChildAllowed::Could not assert child to Node type")
		} else {
			return n.IsChildAllowed(nodeType)
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
    }else if !ok{
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
		e.RemoveNode(node.Id)
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
    if t, ok := TypeMap[nodeType]; !ok {
		return nil, fmt.Errorf("CreateInstanceByType::Type %s is not registered yet.", nodeType)
	} else {
		n := reflect.New(t).Elem().Interface()
		node, ok := n.(Node)

		if !ok {
			return nil, fmt.Errorf("CreateInstanceByType::Unable to create Node for nodetype %s.", nodeType)
		}

		protoSession, protoColl := e.GetMgoSession(PROTOS_SCOPE)
		defer protoSession.Close()

		query := protoColl.Find(bson.M{"tn": nodeType})
		if cnt, err := query.Count(); err != nil {
			return nil, err
		} else if cnt == 0 {
			return nil, fmt.Errorf("CreateInstanceByType::Prototype for NodeType %s is not available", nodeType)
		}

		if err := query.One(&node); err != nil {
			return nil, err
		}

		node.SetId(bson.NewObjectId().Hex())
		return node, nil
    }

    return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNewNode(scope, parentId, name, nodeType string) (Node, error) {
	if node, err := e.CreateInstanceByType(nodeType); err != nil{
        return err
	}else{
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
func (e *Engine) RegisterAllRoutes(scope string, router mux.Router) error {
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

	mainRouter = mux.NewRouter().StrictSlash(false)
	systemRouter:= mux.NewRouter()

	mainRouter.PathPrefix("/nodes").Handler(negroni.New(
		negroni.NewRecovery(),
		//negroni.HandlerFunc(MyMiddleware),
		negroni.NewLogger(),
		negroni.Wrap(systemRouter),
		))

	sysRouter := systemRouter.PathPrefix("/nodes").Subrouter()
	if err := e.RegisterAllRoutes(SYSTEM_SCOPE,sysRouter); err != nil {
		return err
	}

	e.negroni = negroni.Classic()
	e.negroni.UseHandler(mainRouter)
	e.negroni.Run(connection)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	return fn(e)
}

////////////////////////////////////////////////////////////////////////////////
func NewEngine(config *NodesConfig) (*Engine, error) {
	eng := Engine{}

	eng.SetLogger(eng.GetLogger("engine"))

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

	return &eng, nil
}
