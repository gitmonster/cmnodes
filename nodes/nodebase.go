package nodes

import "github.com/denkhaus/cmnodes/render"

type NodeBase struct {
	Id            string `bson:"_id"`
	Parent        string `bson:"p"`
	Order         int    `bson:"o"`
	MimeType      string `bson:"m"`
	TypeName      string `bson:"tn"`
	Route         string `bson:"r"`
	RegisterRoute bool   `bson:"rr"`
	render *render.Render
	engine *Engine
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) assembleRoute() string {
	return n.engine.AssembleRouteFor(n.Id)
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Move(parentId string) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) Delete() error {
	return nil
}