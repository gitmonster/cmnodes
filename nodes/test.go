package nodes

import (
	"path"

	"github.com/gitmonster/cmnodes/helper"
)

func (e *Engine) Test() error {

	set := CriteriaSet{}
	if err := set.LoadFromFile(path.Join(e.StartupDir, "system.toml")); err != nil {
		return err
	}

	if ex, err := e.NodeExists(&set.Criteriae[0]); err != nil {
		return err
	} else {
		helper.Inspect(ex)
	}

	return nil
}
