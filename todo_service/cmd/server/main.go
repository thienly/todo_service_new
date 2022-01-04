package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	grpcPort, _ := strconv.Atoi(os.Getenv("GRPC_PORT"))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db, err := loadDatabase(logger)
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		os.Exit(1)
	}

	go func(logger zerolog.Logger, db *sql.DB, port int) {
		err := startGrpcServer(logger, db, port)
		if err != nil {
			os.Exit(1)
		}
	}(logger, db, grpcPort)

	go func(logger2 zerolog.Logger, grpc, port int) {
		err := startHttpServer(logger2, grpc, port)
		if err != nil {
			os.Exit(1)
		}
	}(logger,grpcPort, httpPort)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Log().Msg("Shutdown application")
	os.Exit(0)
}
