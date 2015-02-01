package nodes

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/denkhaus/validator"
)

type CriteriaSet struct {
	Scope      string     `toml:"Scope" validate:"nonzero"`
	Criteriae  []Criteria `toml:"Criteria"`
	Prototypes []Criteria `toml:"Prototype"`
	Loggable
}

////////////////////////////////////////////////////////////////////////////////
func (c *CriteriaSet) Load(data string) error {
	if _, err := toml.Decode(data, c); err != nil {
		return err
	}

	if err := validator.Validate(c); err != nil {
		return err
	}

	fn := func(cr *Criteria) error {
		if cr.Scope == EMPTY_STRING {
			cr.Scope = c.Scope
		}

		if err := validator.Validate(cr); err != nil {
			return fmt.Errorf("CriteriaSet:Load :: validation error %s :", err)
		}

		return nil
	}

	for idx, p := range c.Prototypes {
		if err := fn(&p); err != nil {
			return err
		}
		c.Prototypes[idx] = p
	}
	for idx, cr := range c.Criteriae {
		if err := fn(&cr); err != nil {
			return err
		}
		c.Criteriae[idx] = cr
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
func (c *CriteriaSet) Ensure(force bool, e *Engine) error {
	fn := func(cr Criteria) error {
		c.Logger.Infof("CriteriaSet:Ensure :: check if node %s already exists.", cr)
		if ex, err := e.NodeExists(&cr); err != nil {
			return err
		} else if !ex {
			c.Logger.Infof("CriteriaSet:Ensure :: node does not exist, create new one.")
			if _, err := e.CreateNode(&cr, false); err != nil {
				return err
			}
		} else {
			c.Logger.Infof("CriteriaSet:Ensure :: node already exists.")
		}

		return nil
	}

	// important:load prototypes first
	for _, p := range c.Prototypes {
		if err := fn(p); err != nil {
			return err
		}
	}
	for _, c := range c.Criteriae {
		if err := fn(c); err != nil {
			return err
		}
	}

	return nil
}
