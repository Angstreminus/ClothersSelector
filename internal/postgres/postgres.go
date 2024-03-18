package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabaseHandler(config *config.Config) (*sqlx.DB, error) {
	ctxTimeout, ctxCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer ctxCancel()
	db, err := sqlx.ConnectContext(ctxTimeout, "postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode))
	if err != nil {
		return nil, err
	}
	return db, nil
}
