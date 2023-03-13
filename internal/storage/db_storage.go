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

var upsertUserSQL = `
	INSERT INTO "user" (login, pwd) 
	VALUES ($1, $2)
	ON CONFLICT (login) DO UPDATE 
		SET pwd = EXCLUDED.pwd
`

func (d *DBStorage) UpsertUser(login string, pwd string) error {
	_, err := d.db.Exec(upsertUserSQL, login, pwd)

	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Healthcheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
