package nodes

import (
    "github.com/gitmonster/cmnodes/render"
    "text"
    "bufio"
)

type NodeBase struct {
	Id            string `bson:"_id"`
	Parent        string `bson:"p"`
	Order         int    `bson:"o"`
	MimeType      string `bson:"m"`
	TypeName      string `bson:"tn"`
	Route         string `bson:"rt"`
	RegisterRoute bool   `bson:"rr"`
	EditRep     Representation   `bson:"er"`
	render        *render.Render
	engine        *Engine
}




////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) RenderEditContent(w *bufio.Writer) error{
    n.EditRep.Render(w)
    return n.EnumerateChilds(func(){

    })
}






////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) assembleRoute() string {
	return n.engine.AssembleRouteFor(n.Id)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Move(parentId string) error {
	return n.engine.MoveNode(Id, parentId)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Delete() error {
	return n.engine.DeleteNode(Id)
}
