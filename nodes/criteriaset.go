package nodes

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type CriteriaSet struct {
	Scope     string     `toml:"Scope"`
	Criterias []Criteria `toml:"Criteria"`
}

////////////////////////////////////////////////////////////////////////////////
func (c *CriteriaSet) Load(data string) error {
	if _, err := toml.Decode(data, c); err != nil {
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
