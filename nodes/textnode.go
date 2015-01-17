package nodes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type TextNode struct {
	NodeBase   `bson:",inline"`
	Content    NodeContent `bson:"c"`
	IsTemplate bool        `bson:"ist"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(TextNode{}, func(engine *Engine) Node {
		node := Node(NewTextNode(engine))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) RegisterRoute(router *mux.Router) {
	router.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		// Assumes you have a template in ./templates called "example.tmpl"
		// $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > templates/example.tmpl
		n.Render.HTML(w, http.StatusOK, "example", nil)
	})
}

////////////////////////////////////////////////////////////////////////////////
func NewTextNode(engine *Engine) *TextNode {
	node := TextNode{}
	node.Init(node, engine)
	node.MimeType = "text/html"
	node.IsTemplate = false

	return &node
}
