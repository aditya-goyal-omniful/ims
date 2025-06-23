package configs

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/migration"
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var DB *postgres.DbCluster

func InitDB(ctx context.Context) {
	dbConfig := postgres.DBConfig{
		Host:                   config.GetString(ctx, "postgres.host"),
		Port:                   config.GetString(ctx, "postgres.port"),
		Username:               config.GetString(ctx, "postgres.user"),
		Password:               config.GetString(ctx, "postgres.password"),
		Dbname:                 config.GetString(ctx, "postgres.name"),
		MaxOpenConnections:     10,
		MaxIdleConnections:     5,
		ConnMaxLifetime:        30 * time.Minute,
		DebugMode:              true, // Set to false in production
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	slaves := []postgres.DBConfig{}

	DB = postgres.InitializeDBInstance(dbConfig, &slaves)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
	)

	migrator, err := migration.InitializeMigrate("file://migrations", dsn)
	if err != nil {
		log.Panic(i18n.Translate(ctx, "Failed to initialize DB migrator: %v"), err)
	}

	if err := migrator.Up(); err != nil {
		log.Panic(i18n.Translate(ctx, "Database migration failed: %v"), err)
	}

	log.Infof(i18n.Translate(ctx, "Connected to database and ran migrations successfully"))
}

func GetDB() *postgres.DbCluster {
	return DB
}