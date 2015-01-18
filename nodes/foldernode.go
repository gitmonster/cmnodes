package nodes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type FolderNode struct {
	NodeBase `bson:",inline"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(FolderNode{}, func(engine *Engine) Node {
		node := Node(NewFolderNode(engine))
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
	node.Init(node, engine)
	node.Name = "FolderNode"
	node.MimeType = "text/html"
	return &node
}
