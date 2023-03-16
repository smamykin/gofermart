package storage

import (
	"context"
	"database/sql"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"time"
)

func NewDBStorage(db *sql.DB) (*DBStorage, error) {

	result := &DBStorage{db: db}
	return result, nil
}

type DBStorage struct {
	db *sql.DB
}

var getUserByLoginSQL = `
	SELECT id, login, pwd
	FROM "user"
	WHERE login = $1
`

func (d *DBStorage) GetUserByLogin(login string) (u entity.User, err error) {
	row := d.db.QueryRow(getUserByLoginSQL, login)
	if row.Err() != nil {
		return u, row.Err()
	}
	err = row.Scan(&u.ID, &u.Login, &u.Pwd)
	if err == nil {

		return u, nil
	}

	if err == sql.ErrNoRows {
		return u, service.ErrUserNotFound
	}

	return u, err
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
