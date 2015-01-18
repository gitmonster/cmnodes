package command

import (
	"github.com/codegangsta/cli"
	"github.com/gitmonster/cmnodes/nodes"
)

func (c *Commander) NewServeCommand() {
	c.Register(cli.Command{
		Name:  "serve",
		Usage: "Start the server",
		Action: func(ctx *cli.Context) {
			c.Execute(func(engine *nodes.Engine) error {
				return engine.Serve()
			}, ctx)
		},
	})
}
