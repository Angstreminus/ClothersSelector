package main

import (
	"fmt"
	"log"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/server"
	"github.com/Angstreminus/ClothersSelector/logger"
)

func main() {
	cfg, err := config.NewConfig()
	fmt.Println(cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}
	zaplog := logger.MustInitLogger()
	zaplog.ZapLogger.Info("Logger initialized")
	Server := server.NewServer(cfg, zaplog)
	Server.MustRun()
}
