package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/db/sql/migration"
	"github.com/omniful/go_commons/db/sql/postgres"
)

var DB *postgres.DbCluster

func InitDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}


	dbConfig := postgres.DBConfig{
		Host:                   os.Getenv("POSTGRES_HOST"),
		Port:                   os.Getenv("POSTGRES_PORT"),
		Username:               os.Getenv("POSTGRES_USER"),
		Password:               os.Getenv("POSTGRES_PASSWORD"),
		Dbname:                 os.Getenv("POSTGRES_DB"),
		MaxOpenConnections:     10,
		MaxIdleConnections:     5,
		ConnMaxLifetime:        30 * time.Minute,
		DebugMode:              true, // Set to false in production
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	slaves := []postgres.DBConfig{} // No read replicas for now

	DB = postgres.InitializeDBInstance(dbConfig, &slaves)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Dbname)

	migrator, err := migration.InitializeMigrate("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to initialize DB migrator: %v", err)
	}

	if err := migrator.Up(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	fmt.Println("Connected to database and ran migrations successfully")
}

func GetDB() *postgres.DbCluster {
	return DB
}
