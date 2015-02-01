package nodes

import (
	"net/http"

	"github.com/gitmonster/cmnodes/helper"
	"github.com/gorilla/mux"
)

type PrototypeNode struct {
	NodeBase `bson:",inline"`
	Template BaseData `bson:"tp"`
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
func (n *PrototypeNode) Apply(crit *Criteria) error {
	if err := helper.BsonTransfer(crit.Template, &n.Template); err != nil {
		return err
	}

	return n.NodeBase.Apply(crit)
}

////////////////////////////////////////////////////////////////////////////////
func (n *PrototypeNode) RegisterRoute(route string, router *mux.Router) {
	router.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
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
