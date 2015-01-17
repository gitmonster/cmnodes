package nodes

import (
	"bufio"
	"time"

	"github.com/gorilla/mux"
)

const (
	EMPTY_STRING = ""
)

type Node interface {
	RenderEditContent(w *bufio.Writer) error
	IsChildAllowed(typeName string) bool
	SetParentId(parentId string)
	Move(parentId string) error
	SetName(name string)
	NewObjectId()
	RegisterRoute(router *mux.Router)
	Remove() error
}

const (
	DURATION_NULL  = time.Second * 0
	DURATION_DAY   = time.Hour * 24
	DURATION_WEEK  = DURATION_DAY * 7
	DURATION_MONTH = DURATION_DAY * 30
)

const (
	SYSTEM_SCOPE = "nodes"
	PROTOS_SCOPE = "protos"
)

// 	acct.Path("/profile").HandlerFunc(ProfileHandler)
//
//
//
// subRouter := e.mux.PathPrefix("/nodes").Subrouter()
//
// subRouter.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "test1") })
// subRouter.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "test2") })

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
