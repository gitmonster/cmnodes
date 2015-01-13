package nodes


import "reflect"

type SCSSStyleNode struct {
	TextNode
}

////////////////////////////////////////////////////////////////////////////////
func init(){
    RegisterNode(new(SCSSStyleNode))
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) SetupRendering() {
	n.engine.mux.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Content here."))
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *TextNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func NewSCSSStyleNode() *SCSSStyleNode {
	node := SCSSStyleNode{}
	node.TypeName = reflect.TypeOf(node).Name()
	node.render = render.New()
	node.MimeType = "text/css"
	node.IsTemplate = false
	node.engine = engine

	return &node
}
