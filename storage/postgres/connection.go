package postgres

import (
	"auth-service/pkg/config"
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func ConnectPostgres(config config.Config) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
