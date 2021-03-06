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
func (n *FolderNode) RegisterRoute(route string, router *mux.Router) {
	router.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("This is a folder."))
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewFolderNode(engine *Engine) *FolderNode {
	node := FolderNode{}
	node.Init(node, engine)
	return &node
}
