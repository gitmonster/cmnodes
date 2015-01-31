package command

import (
	"github.com/codegangsta/cli"
	"github.com/gitmonster/cmnodes/helper"
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
					c.Execute(func(engine *nodes.Engine) error {
						path := ctx.String("path")
						if len(path) > 0 {
							var err error
							path, err = helper.AbsPath(path)
							if err != nil {
								return err
							}
						}
						return engine.ImportSystem(ctx.Bool("force"), path)
					}, ctx)
				},
				Flags: []cli.Flag{
					cli.BoolFlag{"force, f", "Remove any existing system structure before importing.", ""},
					cli.StringFlag{"path, p", "", "Path of the system toml file. When empty the app looks for system.toml in the current directory.", ""},
				},
			},
		},
	})
}
