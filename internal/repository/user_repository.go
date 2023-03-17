package repository

import (
	"database/sql"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *sql.DB
}

var getUserByLoginSQL = `
	SELECT id, login, pwd
	FROM "user"
	WHERE login = $1
`

var upsertUserSQL = `
	INSERT INTO "user" (login, pwd) 
	VALUES ($1, $2)
	ON CONFLICT (login) DO UPDATE 
		SET pwd = EXCLUDED.pwd
	RETURNING id, login, pwd
`

func (d *UserRepository) GetUserByLogin(login string) (u entity.User, err error) {
	row := d.db.QueryRow(getUserByLoginSQL, login)
	if row.Err() != nil {
		return u, row.Err()
	}
	err = row.Scan(&u.ID, &u.Login, &u.Pwd)
	if err == nil {

		return u, nil
	}

	if err == sql.ErrNoRows {
		return u, service.ErrEntityIsNotFound
	}

	return u, err
}

func (d *UserRepository) UpsertUser(login, pwd string) (user entity.User, err error) {
	row := d.db.QueryRow(upsertUserSQL, login, pwd)
	if row.Err() != nil {
		return user, row.Err()
	}
	err = row.Scan(&user.ID, &user.Login, &user.Pwd)
	if err == nil {
		return user, nil
	}

	return user, err
}
