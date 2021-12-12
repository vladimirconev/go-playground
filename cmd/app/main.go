package app

import "github.com/urfave/cli/v2"

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "GOLANG playground"
	app.Commands = []*cli.Command{
		&server,
	}

	return app
}
