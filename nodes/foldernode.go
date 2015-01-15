package nodes

import (
	"net/http"
	"reflect"

	"github.com/gitmonster/cmnodes/render"
)

type FolderNode struct {
	NodeBase
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(new(FolderNode))
}

////////////////////////////////////////////////////////////////////////////////
func (n *FolderNode) IsChildAllowed(typeName string) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////
func (n *FolderNode) SetupRendering() {
	n.engine.mux.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
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
