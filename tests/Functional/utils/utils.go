package utils

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/container"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

func TruncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`TRUNCATE TABLE "user" RESTART IDENTITY CASCADE ;`)
	require.Nil(t, err)
	_, err = db.Exec(`TRUNCATE TABLE "order" RESTART IDENTITY CASCADE;`)
	require.Nil(t, err)
}

var cont = &container.Container{}

func GetContainer(t *testing.T) *container.Container {
	if cont.IsOpen() {
		return cont
	}

	logger := zerolog.Nop()
	newContainer, err := container.NewContainer(&logger)
	require.NoError(t, err)
	t.Cleanup(func() {
		newContainer.Close()
		cont = &container.Container{}
	})
	cont = newContainer

	return newContainer
}

func InsertUser(t *testing.T, db *sql.DB, user entity.User) entity.User {
	if user.Login == "" {
		user.Login = "login" + strconv.Itoa(rand.Int())
	}
	if user.Pwd == "" {
		user.Pwd = "pwd " + strconv.Itoa(rand.Int())
	}

	row := db.QueryRow(`INSERT INTO "user" (login, pwd) VALUES ($1,$2) RETURNING id`, user.Login, user.Pwd)
	require.NoError(t, row.Err())
	err := row.Scan(&user.ID)
	require.NoError(t, err)
	return user
}
