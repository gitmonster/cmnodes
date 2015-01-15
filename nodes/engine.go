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
	"labix.org/v2/bson"
)

////////////////////////////////////////////////////////////////////////////////
type EngineFunc func(engine *Engine) error

type Engine struct {
    Loggable
	MongoSessionProvider
	mux     *mux.Router
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
func (e *Engine) EnumerateChilds(nodeId string, fn EnumFunc) error {
    session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
	defer session.Close()

	iter:= coll.Find(bson.M{"p":nodeId}).Iter()
	if err:= iter.Err(); err != nil{
	    return err
    }

    node := NodeBase{}
    for iter.Next(&node){

    }

    if err:= iter.Close(); err != nil{
	    return err
    }

    return nil
}
////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CheckNodeIsNoChild(nodeId string, parentNodeId string) (bool, error) {
    if nodeId == parentNodeId{
        return false, nil
    }

    session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
	defer session.Close()

	iter:= coll.Find(bson.M{"p":nodeId}).Iter()
	if err:= iter.Err(); err != nil{
	    return err
    }

    node := NodeBase{}
    for iter.Next(&node){
        if nodeId == node.Id{
           return false, nil
        }
    }

    if err:= iter.Close(); err != nil{
	    return err
    }

    return true, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) NodeExists(nodeId string) (bool,error) {
    session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
	defer session.Close()

	if err, cnt:= coll.FindId(srcNodeId).Count(); err != nil{
	    return false, err
	}else if cnt == 1{
	    return true, nil
	}

	return false, nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) MoveNode(srcNodeId, targetNodeId string) error {
    session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
	defer session.Close()

    // check if source node exists
    if ex, err:= e.NodeExists(srcNodeId); err != nil{
        return err
    }else if ! ex{
        return fmt.Errorf("MoveNode::Source node %s doesn't exist.", targetNodeId)
    }

    // check if target node exists
    if ex, err:= e.NodeExists(targetNodeId); err != nil{
        return err
    }else if ! ex{
        return fmt.Errorf("MoveNode::Target node %s doesn't exist.", targetNodeId)
    }

    if ok:= e.CheckNodeIsNoChild(srcNodeId, targetNodeId); err != nil{
        return err
    }else if !ok{
        return fmt.Errorf("MoveNode::Source node %s is already a child of Target node %s.", srcNodeId, targetNodeId)
    }

    set:= bson.M{"$set": bson.M{
        "p": targetNodeId, //set parent to targetNodeId
        "o": 0, //set order to 0
    }}

	if err:= coll.UpdateId(srcNodeId,set); err != nil{
        return err
	}

    return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) RemoveNode(nodeId string) error {
    session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
	defer session.Close()

	// delete childs
	iter:= coll.Find(bson.M{"p":nodeId}).Iter()
	if err:= iter.Err(); err != nil{
	    return err
    }

    node := NodeBase{}
    for iter.Next(&node){
        e.RemoveNode(node.Id)
    }

    if err:= iter.Close(); err != nil{
	    return err
    }

    // delete node
    if err:= coll.RemoveId(nodeId); err != nil{
        return err
    }

    return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNewNode(parentId, name, nodeType string) (Node, error) {
	if t, ok := TypeMap[nodeType]; !ok {
		return nil, fmt.Errorf("CreateNode::Type %s is not registered yet.", nodeType)
	} else {
		n := reflect.New(t).Elem().Interface()
		node, ok := n.(Node)

		if !ok{
		    return nil, fmt.Errorf("CreateNode::Unable to create Node for nodetype %s." nodeType)
		}

		node.SetId(bson.NewObjectId().Hex())
		node.SetParentId(parentId)
		node.SetName(name)

		session, coll := e.GetMgoSession(NODES_COLLECTION_NAME)
		defer session.Close()

		if err := coll.Insert(n); err != nil{
            return nil, err
		}

		return node, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) RegisterAllRoutes() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) AssembleRouteFor(nodeId string) string {
	return EMPTY_STRING
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Startup(connection string) error {
	e.mux = mux.NewRouter()

	mainRouter := mux.NewRouter()
	subRouter := mainRouter.PathPrefix("/").Subrouter()

	subRouter.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "test1") })
	subRouter.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "test2") })

	mainRouter.Handle("/", mainRouter)

	if err := e.RegisterAllRoutes(); err != nil {
		return err
	}

	e.negroni = negroni.Classic()
	e.negroni.UseHandler(e.mux)
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
