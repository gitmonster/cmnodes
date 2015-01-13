package nodes

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type Engine struct {
	mux     *http.ServeMux
	negroni *negroni.Negroni
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) AssembleRouteFor(nodeId string) string {

	return EMPTY_STRING
}

//
// func (e *Engine) RenderData(node Node) error {
//   mux.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
//     r.Data(w, http.StatusOK, []byte("Some binary data here."))
//     })
// 	return nil
// }

//
// mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
// w.Write([]byte("Welcome, visit sub pages now."))
// })
//
//
//
//     mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
//       r.JSON(w, http.StatusOK, map[string]string{"hello": "json"})
//       })
//
//       mux.HandleFunc("/jsonp", func(w http.ResponseWriter, req *http.Request) {
//         r.JSONP(w, http.StatusOK, "callbackName", map[string]string{"hello": "jsonp"})
//         })
//
//         mux.HandleFunc("/xml", func(w http.ResponseWriter, req *http.Request) {
//           r.XML(w, http.StatusOK, ExampleXml{One: "hello", Two: "xml"})
//           })
//
//           mux.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
//             // Assumes you have a template in ./templates called "example.tmpl"
//             // $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > templates/example.tmpl
//             r.HTML(w, http.StatusOK, "example", nil)
//             })
//
//             http.ListenAndServe("0.0.0.0:3000", mux)
//
////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Startup() error {
	e.negroni = negroni.Classic()
	e.negroni.UseHandler(e.mux)
	e.negroni.Run(":3000")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func NewEngine() *Engine {

	eng := &Engine{}
	eng.mux = http.NewServeMux()
	return eng
}
