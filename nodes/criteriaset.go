package nodes

import (
	"github.com/BurntSushi/toml"
)

type CriteriaSet struct {
	Criterias []Criteria
}

func (c *CriteriaSet) Load() error {
	var set CriteriaSet
	if _, err := toml.Decode(blob, &set); err != nil {
		return err
	}
}
