package storage

import (
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	DatabaseName string
	Port         string
	Host         string
	Password     string
	User         string
}

func (p *PostgresConfig) Dialector() gorm.Dialector {
	return postgres.New(
		postgres.Config{
			DSN: fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
				p.Host,
				p.Port,
				p.User,
				p.DatabaseName,
				p.Password,
			),
			DriverName: "postgres",
		},
	)
}

func SetupDatabase(p *PostgresConfig, lg *zap.SugaredLogger) (*gorm.DB, error) {
	return gorm.Open(p.Dialector(), &gorm.Config{
		Logger: logger.New(
			zap.NewStdLog(lg.Desugar()),
			logger.Config{
				LogLevel: logger.Info,
			},
		),
	})
}
