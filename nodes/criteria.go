package nodes

import "labix.org/v2/mgo/bson"

type Criteria struct {
	theMap bson.M
	scope  string
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) Scope() string {
	if c.scope != EMPTY_STRING {
		return c.scope
	}
	panic("Criteria:: Scope is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithId(id string) *Criteria {
	c.theMap["_id"] = id
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) Id() string {
	if ob, ok := c.theMap["_id"]; ok {
		return ob.(string)
	}
	panic("Criteria:: Id is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithName(name string) *Criteria {
	c.theMap["nm"] = name
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) Name() string {
	if ob, ok := c.theMap["nm"]; ok {
		return ob.(string)
	}
	panic("Criteria:: Name is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithParentId(parentId string) *Criteria {
	c.theMap["p"] = parentId
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) ParentId() string {
	if ob, ok := c.theMap["p"]; ok {
		return ob.(string)
	}
	panic("Criteria:: ParentId is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithOrder(order int) *Criteria {
	c.theMap["o"] = order
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) Order() int {
	if ob, ok := c.theMap["o"]; ok {
		return ob.(int)
	}
	panic("Criteria:: Order is not yet defined!")
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) WithNodeType(nodeType string) *Criteria {
	c.theMap["tn"] = nodeType
	return c
}

////////////////////////////////////////////////////////////////////////////////
func (c *Criteria) NodeType() string {
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
	return &Criteria{scope: scope}
}
