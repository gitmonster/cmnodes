package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/denkhaus/tcgl/applog"
	"github.com/gitmonster/cmnodes/command"
	"github.com/gitmonster/cmnodes/nodes"
)

var releaseVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "cmnodes"
	app.Version = releaseVersion
	app.Usage = "A cms and static file site generator powerd by golang."
	app.Flags = []cli.Flag{
		//cli.StringFlag{"group, g", "", "group or container to restrict the command to"},
		//cli.StringFlag{"manifest, m", "", "path to a manifest (.json, .yml, .yaml) file to read from"},
		cli.BoolFlag{"debug, d", "print debug output", ""},
		//cli.BoolFlag{"test, t", "perform selftests"},
		//cli.StringSliceFlag{"peers, C", &cli.StringSlice{}, "a comma-delimited list of machine addresses in the cluster (default: {\"127.0.0.1:4001\"})"},
	}

	if cnf, err := nodes.CreateNodesConfig(); err != nil {
		applog.Errorf("Configuration error:: load config %s", err.Error())
		return
	} else {
		if cmdr, err := command.NewCommander(app, cnf); err != nil {
			applog.Errorf("Startup error:: %s", err.Error())
			return
		} else {
			cmdr.Run(os.Args)
		}
	}
}
