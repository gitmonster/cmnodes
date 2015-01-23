package nodes

import "labix.org/v2/mgo/bson"

type Criteria struct {
	BaseData `toml:"-"`
	theMap   bson.M `toml:"-"`
	Data BaseData `toml:"Data"` //Payload for PrototypeNodes
}


////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) HasProtoData() bool {
	return c.Data != nil
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithScope(scope string) *Criteria {
	c.Scope = scope
	c.theMap["sp"] = scope
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) HasScope() bool {
	_, ok := c.theMap["sp"]
	return ok
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetScope() string {
	if ob, ok := c.theMap["sp"]; ok {
		return ob.(string)
	}

	panic("Criteria:: Scope is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithId(id string) *Criteria {
	c.Id = id
	c.theMap["_id"] = id
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetId() string {
	if ob, ok := c.theMap["_id"]; ok {
		return ob.(string)
	}
	panic("Criteria:: Id is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithName(name string) *Criteria {
	c.Name = name
	c.theMap["nm"] = name
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetName() string {
	if ob, ok := c.theMap["nm"]; ok {
		return ob.(string)
	}
	panic("Criteria:: Name is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithParentId(parentId string) *Criteria {
	c.ParentId = parentId
	c.theMap["p"] = parentId
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetParentId() string {
	if ob, ok := c.theMap["p"]; ok {
		return ob.(string)
	}
	panic("Criteria:: ParentId is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithOrder(order int) *Criteria {
	c.Order = order
	c.theMap["o"] = order
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetOrder() int {
	if ob, ok := c.theMap["o"]; ok {
		return ob.(int)
	}
	panic("Criteria:: Order is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithNodeType(nodeType string) *Criteria {
	c.TypeName = nodeType
	c.theMap["tn"] = nodeType
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetNodeType() string {
	if ob, ok := c.theMap["tn"]; ok {
		return ob.(string)
	}
	panic("Criteria:: NodeType is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) GetSelector() bson.M {
	return c.theMap
}

////////////////////////////////////////////////////////////////////////////////
func NewCriteria(scope string) *Criteria {
	cr := Criteria{}
	cr.Scope = scope
	cr.theMap = bson.M{}
	return &cr
}
