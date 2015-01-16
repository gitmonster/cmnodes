package nodes

import (
	"github.com/gitmonster/cmnodes/render"

	"bufio"
)

type NodeBase struct {
	Id            string         `bson:"_id"`
	Parent        string         `bson:"p"`
	Order         int            `bson:"o"`
	MimeType      string         `bson:"m"`
	TypeName      string         `bson:"tn"`
	Route         string         `bson:"rt"`
	RegisterRoute bool           `bson:"rr"`
	EditRep       Representation `bson:"er"`
	Scope         string         `bson:"sp"`
	render        *render.Render
	engine        *Engine
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) RenderEditContent(w *bufio.Writer) error {
	if err := n.EditRep.Render(w); err != nil {
		return err
	}
	return n.EnumerateChilds(func(node Node) error {
		return node.RenderEditContent(w)
	})
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) EnumerateChilds(fn EnumFunc) error {
	return n.engine.EnumerateChilds(n.Scope, n.Id, fn)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) assembleRoute() string {
	return n.engine.AssembleRouteFor(n.Scope, n.Id)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Move(parentId string) error {
	return n.engine.MoveNode(n.Scope, n.TypeName, n.Id, parentId)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Remove() error {
	return n.engine.RemoveNode(n.Scope, n.Id)
}
