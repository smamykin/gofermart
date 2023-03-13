package storage

import (
	"context"
	"database/sql"
	"time"
)

func NewDBStorage(db *sql.DB) (*DBStorage, error) {

	result := &DBStorage{db: db}
	return result, nil
}

type DBStorage struct {
	db *sql.DB
}

func (d *DBStorage) Healthcheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
