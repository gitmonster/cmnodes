package nodes

import (
	"net/http"
	"github.com/gorilla/mux"
)

type FileNode struct {
	NodeBase   `bson:",inline"`
	ResourceId string `bson:"rid"`
    RefNodes   []string `bson:"refs"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(FileNode{}, func(engine *Engine) Node {
		node := Node(NewFileNode(engine))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *FileNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func (n *FileNode) AddReference(refNodeId string) error {
    return n.Engine.AddReference(n.Scope, n.Id, refNodeId)
}

////////////////////////////////////////////////////////////////////////////////
func (n *FileNode) RegisterRoute(route string, router *mux.Router) {
	router.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		// Assumes you have a template in ./templates called "example.tmpl"
		// $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > templates/example.tmpl
		n.Render.HTML(w, http.StatusOK, "example", nil)
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewFileNode(engine *Engine) *FileNode {
	node := FileNode{}
	node.Init(node, engine)
	return &node
}