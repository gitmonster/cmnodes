package nodes

import (
	"net/http"
	"reflect"

	"github.com/gitmonster/cmnodes/render"
	"github.com/gorilla/mux"
)

type FolderNode struct {
	NodeBase
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(FolderNode{}, func() Node {
		node := Node(new(FolderNode))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *FolderNode) IsChildAllowed(typeName string) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////
func (n *FolderNode) RegisterRoute(router *mux.Router) {
	router.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("This is a folder."))
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewFolderNode(engine *Engine) *FolderNode {
	node := FolderNode{}
	node.TypeName = reflect.TypeOf(node).Name()
	node.render = render.New()
	node.MimeType = "text/html"
	node.engine = engine

	return &node
}
