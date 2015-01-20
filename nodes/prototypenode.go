package nodes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type PrototypeNode struct {
	NodeBase `bson:",inline"`
	Proto    BaseData `bson:"pr" toml:"Proto"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(PrototypeNode{}, func(engine *Engine) Node {
		node := Node(NewPrototypeNode(engine))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *PrototypeNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func (n *PrototypeNode) RegisterRoute(router *mux.Router) {
	router.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		// Assumes you have a template in ./templates called "example.tmpl"
		// $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > templates/example.tmpl
		n.Render.HTML(w, http.StatusOK, "example", nil)
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewPrototypeNode(engine *Engine) *PrototypeNode {
	node := PrototypeNode{}
	node.Init(node, engine)
	return &node
}
