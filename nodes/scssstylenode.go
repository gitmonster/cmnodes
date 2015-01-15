package nodes

import (
	"net/http"
	"reflect"

	"github.com/gitmonster/cmnodes/render"
)

type SCSSStyleNode struct {
	TextNode
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(new(SCSSStyleNode))
}

////////////////////////////////////////////////////////////////////////////////
func (n *SCSSStyleNode) RegisterRoute() {
	n.engine.mux.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Content here."))
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *SCSSStyleNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func NewSCSSStyleNode(engine *Engine) *SCSSStyleNode {
	node := SCSSStyleNode{}
	node.TypeName = reflect.TypeOf(node).Name()
	node.render = render.New()
	node.MimeType = "text/css"
	node.IsTemplate = false
	node.engine = engine

	return &node
}
