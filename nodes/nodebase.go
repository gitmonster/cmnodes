package nodes

import (
	"bufio"

	"github.com/gitmonster/cmnodes/render"
	"labix.org/v2/mgo/bson"
)

type NodeBase struct {
	BaseData `bson:",inline"`
	Render   *render.Render `bson:"-" toml:"-"`
	Engine   *Engine        `bson:"-" toml:"-"`
	Loggable `bson:"-" toml:"-"`
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Init(inst interface{}, e *Engine) {
	n.TypeName = GetTypeName(inst)
	n.SetLogger(e.NewLogger(n.TypeName))
	n.Render = render.New()
	n.Engine = e
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Apply(crit *Criteria) error {
	if crit.HasObjectId() {
		n.Id = crit.GetObjectId()
	}
	if crit.HasScope() {
		n.Scope = crit.GetScope()
	}
	if crit.HasName() {
		n.Name = crit.GetName()
	}
	if crit.HasNodeType() {
		n.TypeName = crit.GetNodeType()
	}
	if crit.HasOrder() {
		n.Order = crit.GetOrder()
	}

	return nil
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
func (n *NodeBase) GetObjectId() string {
	return n.Id
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) SetOrder(order int) {
	n.Order = order
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
	n.ParentId = parentId
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) GetParentId() string {
	return n.ParentId
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) MustRegisterRoute() bool {
	return n.RegRoute
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
func (n *NodeBase) Move(parentId string) error {
	return n.Engine.MoveNode(n.Scope, n.TypeName, n.Id, parentId)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Remove() error {
	return n.Engine.RemoveNode(n.Scope, n.Id)
}
