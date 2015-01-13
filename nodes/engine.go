package nodes

import (
	"net/http"
	"fmt"
    "reflect"
	"github.com/codegangsta/negroni"
)



type Engine struct {
	mux     *http.ServeMux
	negroni *negroni.Negroni
}

var( TypeMap map[string]reflect.Type)


////////////////////////////////////////////////////////////////////////////////
func init(){
    TypeMap = make(map[string]reflect.Type)
}

////////////////////////////////////////////////////////////////////////////////
func RegisterNode(node interface{}){
    t:=reflect.TypeOf(node)
    TypeMap[t.Name()] = t
}


////////////////////////////////////////////////////////////////////////////////
func (e *Engine) CreateNode(nodeType string) (interface{}, error) {
    if t, ok:= TypeMap[nodeType]; !ok{
        return nil, fmt.Errorf("Type %s is not registered yet.", nodeType)
    }
    i := reflect.New(yourtype).Elem().Interface()
	return i , nil
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
func (e *Engine) Startup() error {
	if err:= e.RegisterAllRoutes(); err != nil{
	    return err
	}

	e.negroni = negroni.Classic()
	e.negroni.UseHandler(e.mux)
	e.negroni.Run(":3000")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func NewEngine() *Engine {
	eng := Engine{}
	eng.mux = http.NewServeMux()
	return &eng
}
