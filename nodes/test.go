package nodes

import "path"

func (e *Engine) Test() error {

	set := CriteriaSet{}
	if err := set.LoadFromFile(path.Join(e.StartupDir, "system.toml")); err != nil {
		return err
	}
	Inspect(set)

	return nil
}
