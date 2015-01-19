package nodes

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/validator.v2"
)

type CriteriaSet struct {
	Scope      string     `toml:"Scope" validate:"nonzero"`
	Criterias  []Criteria `toml:"Criteria"`
	Prototypes []Criteria `toml:"Prototype"`
}

////////////////////////////////////////////////////////////////////////////////
func (c *CriteriaSet) Load(data string) error {
	if _, err := toml.Decode(data, c); err != nil {
		return err
	}
	if err := validator.Validate(c); err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (c *CriteriaSet) LoadFromFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	if buf, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return c.Load(string(buf))
	}
}

////////////////////////////////////////////////////////////////////////////////
func (c *CriteriaSet) Ensure(e *Engine) error {
	// load prototypes first
	for _, p := range c.Prototypes {
		if !p.HasScope() {
			p.WithScope(c.Scope)
		}
		if err := validator.Validate(p); err != nil {
			return err
		}
		if ex, err := e.NodeExists(&p); err != nil {
			return nil
		} else if !ex {
			if _, err := e.CreateNode(&p, false); err != nil {
				return err
			}
		}
	}
	for _, c := range c.Criterias {
		if !c.HasScope() {
			c.WithScope(c.Scope)
		}
		if err := validator.Validate(c); err != nil {
			return err
		}
		if ex, err := e.NodeExists(&c); err != nil {
			return nil
		} else if !ex {
			// abort if prototype does not exist.
			if _, err := e.CreateNode(&c, true); err != nil {
				return err
			}
		}
	}

	return nil
}
