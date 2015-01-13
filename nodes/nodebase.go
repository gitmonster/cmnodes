package nodes

import "github.com/denkhaus/cmnodes/render"

type NodeBase struct {
	Id            string `bson:"_id"`
	Parent        string `bson:"p"`
	Order         int    `bson:"o"`
	MimeType      string `bson:"m"`
	RegisterRoute bool   `bson:"rr"`
	Route         string `bson:"r"`

	render *render.Render
	engine *Engine
}

////////////////////////////////////////////////////////////////////////////////
func (n *NodeBase) assembleRoute() string {
	return n.engine.AssembleRouteFor(n.Id)
}
