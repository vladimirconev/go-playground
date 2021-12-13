package app

import (
	"fmt"

	"example.com/playground/pkg/rest"
	"example.com/playground/pkg/storage"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var server = cli.Command{
	Name: "server",
	Action: func(c *cli.Context) error {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}

		lg := logger.Sugar()

		ctx := c.Context

		postgresConfig := &storage.PostgresConfig{
			DatabaseName: c.String("postgres-db"),
			Host:         c.String("postgres-host"),
			Port:         c.String("postgres-port"),
			User:         c.String("postgres-user"),
			Password:     c.String("postgres-password"),
		}

		db, err := storage.SetupDatabase(ctx, postgresConfig, lg)
		if err != nil {
			return err
		}

		r := rest.SetupRouteHandlers(&rest.RouteHandlers{
			CreateOffer: storage.NewCreateOfferService(db),
			UpdateOffer: storage.NewUpdateOfferService(db),
			GetOffer:    storage.NewGetOfferService(db),
			DeleteOffer: storage.NewDeleteOfferService(db),
		}, lg)

		return r.Run(fmt.Sprintf(":%s", c.String("server-port")))
	},
	Flags: []cli.Flag{
		&cli.StringFlag{EnvVars: []string{"SERVER_PORT"}, Name: "server-port", Value: "3456"},
		&cli.StringFlag{EnvVars: []string{"POSTGRES_HOST"}, Name: "postgres-host", Value: "localhost"},
		&cli.StringFlag{EnvVars: []string{"POSTGRES_PORT"}, Name: "postgres-port", Value: "5432"},
		&cli.StringFlag{EnvVars: []string{"POSTGRES_DB"}, Name: "postgres-db", Value: "offers_db"},
		&cli.StringFlag{EnvVars: []string{"POSTGRES_USER"}, Name: "postgres-user", Value: "postgres"},
		&cli.StringFlag{EnvVars: []string{"POSTGRES_PASSWORD"}, Name: "postgres-password", Value: "your_password"},
	},
}
