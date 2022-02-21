package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var (
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var httpPort, grpcPort int64
	flag.Int64Var(&httpPort, "HTTP_PORT", 8088, "http port")
	flag.Int64Var(&grpcPort, "GRPC_PORT", 8087, "grpc port")
	flag.Parse()

	db, err := loadDatabase(logger)
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		fmt.Errorf("error while starting application: %v", err)
		os.Exit(1)
	}

	go func(logger zerolog.Logger, db *sqlx.DB, port int) {
		err := startGrpcServer(logger, db, port)
		if err != nil {
			os.Exit(1)
		}
	}(logger, db, int(grpcPort))

	go func(logger2 zerolog.Logger, grpc, port int) {
		err := startHttpServer(logger2, grpc, port)
		if err != nil {
			os.Exit(1)
		}
	}(logger, int(grpcPort), int(httpPort))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Log().Msg("Shutdown application")
	os.Exit(0)
}
