package db

import (
	"fmt"
	"log"
	"time"
	"truckx/internal/models"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func Connect(cfg models.Config) (*sqlx.DB, error) {
	// allow some time for pg server to boot
	time.Sleep(5 * time.Second)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)
	log.Println("dns", dsn)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
