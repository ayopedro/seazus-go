package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ayopedro/seazus-go/internal/config"
	_ "github.com/lib/pq"
)

func New(cfg config.DBConfig) (*sql.DB, error) {
	dsn := cfg.URI
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
