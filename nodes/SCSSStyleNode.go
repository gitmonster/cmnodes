package nodes

type SCSSStyleNode struct {
	TextNode
}

////////////////////////////////////////////////////////////////////////////////
func (n *SCSSStyleNode) Render() error {

	n.MimeType = "text/css"
	n.IsTemplate = false
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func NewSCSSStyleNode() *SCSSStyleNode {
	node := &SCSSStyleNode{}
	return node
}
