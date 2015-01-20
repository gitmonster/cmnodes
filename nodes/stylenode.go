package nodes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type StyleNode struct {
	TextNode `bson:",inline"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(StyleNode{}, func(engine *Engine) Node {
		node := Node(NewStyleNode(engine))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *StyleNode) RegisterRoute(router *mux.Router) {
	router.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Content here."))
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *StyleNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func NewStyleNode(engine *Engine) *StyleNode {
	node := StyleNode{}
	node.Init(node, engine)

	return &node
}
