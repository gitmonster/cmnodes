package nodes

import (
	"net/http"
	"reflect"

	"github.com/gitmonster/cmnodes/render"
)

type TextNode struct {
	NodeBase
	Content    NodeContent `bson:"c"`
	IsTemplate bool        `bson:"ist"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNode(new(TextNode))
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) SetupRendering() {
	n.engine.mux.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		// Assumes you have a template in ./templates called "example.tmpl"
		// $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > templates/example.tmpl
		n.render.HTML(w, http.StatusOK, "example", nil)
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewTextNode(engine *Engine) *TextNode {
	node := TextNode{}
	node.TypeName = reflect.TypeOf(node).Name()
	node.render = render.New()
	node.MimeType = "text/html"
	node.IsTemplate = false
	node.engine = engine

	return &node
}
