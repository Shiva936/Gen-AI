package main

import (
	"api-server/assets"
	"api-server/config"
	"api-server/handlers"
	datapipeline "api-server/services/data-pipeline"
	"api-server/services/documents"
	"database/sql"
	"embed"
	"log"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func main() {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("unable to open config file.")
	}

	var conf *config.Config
	config.ParseJSON(file, &conf)
	config.Set(conf)

	runMigrations()

	// Mocking Upload Document Flow

	err = documents.NewDocument().SaveAndTriggerDataPipeline(assets.FinanceBenchDataset)
	if err != nil {
		log.Fatal("unable to upload document")
	}

	datapipeline.NewDataPipeline().Start()

	handlers.GetRouter()
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func runMigrations() {

	migrationUserDbUrl := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(migrationUserDbUrl) == "" {
		log.Fatal("DATABASE_URL is not provided")
	}
	db, err := sql.Open("postgres", migrationUserDbUrl)
	if err != nil {
		log.Fatal("PG DB Connection Failed", zap.Error(err))
	}
	// setup database
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Setting Goose Postgres Dialect Failed", zap.Error(err))
	}
	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		log.Fatal("Goose Up Failed", zap.Error(err))
	}

	log.Println("Successfully completed database migrations")
}
