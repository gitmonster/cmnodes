package nodes

import (
	"fmt"

	"github.com/gitmonster/cmnodes/helper"

	"labix.org/v2/mgo/bson"
)

type Criteria struct {
	BaseData `bson:",inline" toml:"Base"`
	//Payload for PrototypeNodes
	Template BaseData `bson:"tp,omitempty" toml:"Template" validate:"-"`
}

//////////////////////////////////////////////////////////////////////////////////
func (c Criteria) String() string {
	return fmt.Sprintf("Criteria with Id '%s' and NodeType '%s' and Name '%s'", c.Id, c.NodeType, c.Name)
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithProtoNodeType(nodeType string) *Criteria {
	if !c.HasTemplate() {
		c.Template = BaseData{NodeType: nodeType}
	}
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) HasTemplate() bool {
	return c.Template != BaseData{}
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithScope(scope string) *Criteria {
	c.Scope = scope
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithMimeType(mimeType string) *Criteria {
	c.MimeType = mimeType
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithNewObjectId() *Criteria {
	c.Id = bson.NewObjectId().Hex()
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithObjectId(id string) *Criteria {
	c.Id = id
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithName(name string) *Criteria {
	c.Name = name
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithParentId(parentId string) *Criteria {
	c.ParentId = parentId
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithOrder(order int) *Criteria {
	c.Order = order
	return c
}

//////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithNodeType(nodeType string) *Criteria {
	c.NodeType = nodeType
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetSelector() bson.M {
	sel := bson.M{}
	if err := helper.BsonTransfer(*c, &sel); err != nil {
		panic(fmt.Sprintf("Criteria:GetSelector :: %s", err))
	}

	return sel
}

////////////////////////////////////////////////////////////////////////////////
func NewCriteria(scope string) *Criteria {
	cr := &Criteria{}
	return cr.WithScope(scope)
}
