package command

import (
	"github.com/codegangsta/cli"
	"github.com/gitmonster/cmnodes/nodes"
)

type Commander struct {
	nodes.Loggable
	engine *nodes.Engine
	config *nodes.NodesConfig
	app    *cli.App
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// func (c *Commander) Execute(fn engine.EngineFunc, ctx *cli.Context) {
// 	//if ctx.GlobalBool("debug") {
// 	//	c.Logger.SetLevel(engine.LevelDebug)
// 	//} else {
// 	//	c.Logger.SetLevel(engine.LevelInfo)
// 	//}
//
// 	if err := c.engine.Execute(fn); err != nil {
// 		c.Logger.Errorf("Execution error:: %s", err.Error())
// 	}
// }

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func NewCommander(app *cli.App, cnf *nodes.NodesConfig) (*Commander, error) {
	cmd := &Commander{app: app, config: cnf}
	if engine, err := nodes.NewEngine(cnf); err != nil {
		return nil, err
	} else {
		cmd.engine = engine
	}

	// cmd.NewStatusCommand()
	// cmd.NewTestCommand()
	// cmd.NewDelegateStatusCommand()
	// cmd.NewScanBlockchainCommand()
	return cmd, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Register(cmd cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Run(args []string) {
	//gosignal.ObserveInterrupt().Then(func() {
	//	logger.Infof("Termination requested - shutting down...")
	//	c.engine.EnableRunning(false)
	//})
	c.app.Run(args)
	c.engine.Session.Close()
}
