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
)

////////////////////////////////////////////////////////////////////////////////
type EngineFunc func(engine *Engine) error

type Engine struct {
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
func RegisterNode(node interface{}) {
	t := reflect.TypeOf(node)
	TypeMap[t.Name()] = t
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNode(nodeType string) (interface{}, error) {
	if t, ok := TypeMap[nodeType]; !ok {
		return nil, fmt.Errorf("Type %s is not registered yet.", nodeType)
	} else {
		i := reflect.New(t).Elem().Interface()
		return i, nil
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
			return nil, fmt.Errorf("Init error:: Mongo Session could not be initialized. :: %s", err.Error())
		} else {
			session.SetMode(mgo.Monotonic, true)
			eng.Session = session
		}
	}

	return &eng, nil
}
