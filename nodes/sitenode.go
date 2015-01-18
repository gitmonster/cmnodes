package nodes

import (
	"github.com/gorilla/mux"
	//"net/http"
)

type SiteNode struct {
	NodeBase `bson:",inline"`
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	RegisterNodeType(SiteNode{}, func(engine *Engine) Node {
		node := Node(NewSiteNode(engine))
		return node
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *SiteNode) RegisterRoute(router *mux.Router) {
	//router.HandleFunc(n.assembleRoute(), func(w http.ResponseWriter, req *http.Request) {
	//	w.Write([]byte("Content here."))
	//})
}

////////////////////////////////////////////////////////////////////////////////
func (n *SiteNode) IsChildAllowed(typeName string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
func NewSiteNode(engine *Engine) *SiteNode {
	node := SiteNode{}
	node.Init(node, engine)
	return &node
}
