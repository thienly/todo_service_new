package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"new_todo_project/pkg/config"

	"github.com/rs/zerolog"
)

func loadDatabase(logger zerolog.Logger) (*sqlx.DB, error) {
	_, err := config.LoadFromJsonOrPanic("config.json")
	if err != nil {
		logger.Fatal().Msg("Can not read config.json file!!")
	}
	db, err := sqlx.Connect("postgres", config.AppCfg.Database.ConnectionString)

	if err != nil {
		logger.Fatal().Msg("Can not open database connection")
		return nil, err
	}
	return db, nil
}
