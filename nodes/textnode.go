package nodes

import (
	"net/http"

	"github.com/denkhaus/cmnodes/render"
)

type TextNode struct {
	NodeBase

	Content    NodeContent `bson:"c"`
	IsTemplate bool        `bson:"ist"`
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) Render() error {
	n.MimeType = "text/html"
	return nil
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
	node := &TextNode{}
	node.engine = engine
	node.render = render.New()

	return node
}
