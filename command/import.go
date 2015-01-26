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
				Name:  "system",
				Usage: "Import system node structure with object prototypes",
				Action: func(ctx *cli.Context) {
					var path string 
					if len(ctx.Args()) > 0 {
    					path = AbsPath(ctx.Args()[0])
  					}
					c.Execute(func(engine *nodes.Engine) error {
						return nil engine.ImportSystem(ctx.Bool("force"),path)
					}, ctx)
				},
				Flags: []cli.Flag{
					cli.BoolFlag{"force, f", "Remove any existing system structure before importing.", ""},
					cli.BoolFlag{"path, p", "Path of the system toml file. When empty the app looks for system.toml in the current directory.", ""},
				},
			},
		},
	})
}
