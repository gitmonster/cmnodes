package command

import (
	"github.com/codegangsta/cli"
	"github.com/gitmonster/cmnodes/nodes"
)

func (c *Commander) NewTestCommand() {
	c.Register(cli.Command{
		Name:  "test",
		Usage: "Run seveeral tests",
		Action: func(ctx *cli.Context) {
			c.Execute(func(engine *nodes.Engine) error {
				return engine.Test()
			}, ctx)
		},
	})
}
