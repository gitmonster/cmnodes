package command

import (
	"github.com/codegangsta/cli"
	"github.com/gitmonster/cmnodes/nodes"
)

func (c *Commander) NewInitProtosCommand() {
	c.Register(cli.Command{
		Name:  "import",
		Usage: "Import command with subcommands",
		Subcommands: []cli.Command{
			{
				Name:  "prototypes",
				Usage: "Import prototypes",
				Action: func(ctx *cli.Context) {
					c.Execute(func(engine *nodes.Engine) error {
						return nil //engine.ImportPrototypes(ctx.Bool("force"))
					}, ctx)
				},
				Flags: []cli.Flag{
					cli.BoolFlag{"force, f", "delete all prototypes in database before importing", ""},
				},
			},
		},
	})
}
