package nodes

import (
	"bufio"

	"github.com/gitmonster/cmnodes/render"
	"labix.org/v2/mgo/bson"
)

type NodeBase struct {
	Id            string         `bson:"_id"`
	Parent        string         `bson:"p"`
	Name          string         `bson:"nm"`
	Order         int            `bson:"o"`
	MimeType      string         `bson:"m"`
	TypeName      string         `bson:"tn"`
	Route         string         `bson:"rt"`
	RegisterRoute bool           `bson:"rr"`
	EditRep       Representation `bson:"er"`
	Scope         string         `bson:"sp"`
	Render        *render.Render `bson:"-"`
	Engine        *Engine        `bson:"-"`
	Loggable      `bson:"-"`
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Init(inst interface{}, e *Engine) {
	n.TypeName = GetNodeTypeName(inst)
	n.SetLogger(e.NewLogger(n.TypeName))
	n.Render = render.New()
	n.Engine = e
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) SetEditTemplate(content string) {
	n.EditRep.Content = content
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) SetObjectId(objectId string) {
	n.Id = objectId
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) NewObjectId() {
	n.Id = bson.NewObjectId().Hex()
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) SetName(name string) {
	n.Name = name
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) SetParentId(parentId string) {
	n.Parent = parentId
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
	return n.Engine.EnumerateChilds(n.Scope, n.Id, fn)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) assembleRoute() string {
	return n.Engine.AssembleRouteFor(n.Scope, n.Id)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Move(parentId string) error {
	return n.Engine.MoveNode(n.Scope, n.TypeName, n.Id, parentId)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Remove() error {
	return n.Engine.RemoveNode(n.Scope, n.Id)
}
