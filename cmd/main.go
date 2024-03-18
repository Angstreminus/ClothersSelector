package main

import (
	"context"
	"log"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/postgres"
	"github.com/Angstreminus/ClothersSelector/internal/server"
	"github.com/Angstreminus/ClothersSelector/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	zaplog := logger.MustInitLogger()
	zaplog.ZapLogger.Info("Logger initialized")
	db, err := postgres.NewDatabaseHandler(cfg)
	if err != nil {
		zaplog.ZapLogger.Error("Error to connect postgres")
	}

	err = db.PingContext(context.Background())
	if err != nil {
		zaplog.ZapLogger.Fatal("Fail to connect")
		log.Fatal(err)
	}
	zaplog.ZapLogger.Info("Connected and ping success")
	Server := server.NewServer(cfg, zaplog)
	Server.MustRun()
}
