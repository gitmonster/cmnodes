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
	GetParentId() string
	Move(parentId string) error
	SetName(name string)
	SetOrder(order int)
	SetEditTemplate(content string)
	SetObjectId(objectId string)
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
)

const (
	OBJECTID_SYSTEM_SITE       = "54bc1c73618ccf2345600005"
	OBJECTID_SYSTEM_PROTOTYPES = "54bc1c3456cdf458cc000453"
	OBJECTID_SYSTEM_TEMPLATES  = "54bc1c73618cc458cc0567f5"
	OBJECTID_SYSTEM_CONTENT    = "54bc1c73618cfc345c00fc34"
)

var (
	NODETYPE_SITE   = GetTypeName(SiteNode{})
	NODETYPE_TEXT   = GetTypeName(TextNode{})
	NODETYPE_STYLE  = GetTypeName(StyleNode{})
	NODETYPE_FOLDER = GetTypeName(FolderNode{})
)
var (
	CRITERIA_SYSTEM_SITE       = NewCriteria(SYSTEM_SCOPE).WithName("System").WithNodeType(NODETYPE_SITE).WithId(OBJECTID_SYSTEM_SITE).WithParentId(EMPTY_STRING).WithOrder(99)
	CRITERIA_SYSTEM_CONTENT    = NewCriteria(SYSTEM_SCOPE).WithId(OBJECTID_SYSTEM_CONTENT).WithName("Content").WithNodeType(NODETYPE_FOLDER).WithParentId(OBJECTID_SYSTEM_SITE).WithOrder(0)
	CRITERIA_SYSTEM_TEMPLATES  = NewCriteria(SYSTEM_SCOPE).WithId(OBJECTID_SYSTEM_TEMPLATES).WithName("Templates").WithNodeType(NODETYPE_FOLDER).WithParentId(OBJECTID_SYSTEM_SITE).WithOrder(1)
	CRITERIA_SYSTEM_PROTOTYPES = NewCriteria(SYSTEM_SCOPE).WithId(OBJECTID_SYSTEM_PROTOTYPES).WithName("Prototypes").WithNodeType(NODETYPE_FOLDER).WithParentId(OBJECTID_SYSTEM_SITE).WithOrder(2)
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
