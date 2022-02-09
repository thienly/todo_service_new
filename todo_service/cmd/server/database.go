package main

import (
	"database/sql"
	"new_todo_project/pkg/config"

	"github.com/rs/zerolog"
)

func loadDatabase(logger zerolog.Logger) (*sql.DB, error) {
	_, err := config.LoadFromJsonOrPanic("config.json")
	if err != nil {
		logger.Fatal().Msg("Can not read config.json file!!")
	}
	db, err := sql.Open("postgres", config.AppCfg.Database.ConnectionString)
	err = db.Ping()
	if err != nil {
		logger.Fatal().Msg("Can not open database connection")
		return nil, err
	}
	return db, nil
}
