package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDBStorage_UpsertUser(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)
	store := c.Storage()

	err := store.UpsertUser("cheesecake", "pwd")
	require.NoError(t, err)
	assertUsersInDB(t, db, []entity.User{
		{ID: 1, Login: "cheesecake", Pwd: "pwd"},
	})

	err = store.UpsertUser("cheesecake", "pwd2")
	require.NoError(t, err)
	assertUsersInDB(t, db, []entity.User{
		{ID: 1, Login: "cheesecake", Pwd: "pwd2"},
	})
	err = store.UpsertUser("cheesecake2", "pwd")
	require.NoError(t, err)
	assertUsersInDB(t, db, []entity.User{
		{ID: 1, Login: "cheesecake", Pwd: "pwd2"},
		{ID: 3, Login: "cheesecake2", Pwd: "pwd"},
	})
}

func TestDBStorage_GetUserByLogin(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	expected := entity.User{
		ID:    1,
		Login: "foo",
		Pwd:   "bar",
	}
	insertUser(t, db, expected)

	store := c.Storage()

	actual, err := store.GetUserByLogin("foo")
	require.Nil(t, err)
	assert.Equal(t, expected, actual)

	_, err = store.GetUserByLogin("baz")
	assert.Equal(t, service.ErrUserIsNotFound, err)

}

func insertUser(t *testing.T, db *sql.DB, expected entity.User) {
	_, err := db.Exec(`INSERT INTO "user" (login, pwd) VALUES ($1,$2);`, expected.Login, expected.Pwd)
	require.Nil(t, err)
}

func assertUsersInDB(t *testing.T, db *sql.DB, expected []entity.User) {
	getUsersSQL := `
		SELECT id, login, pwd
		FROM public."user"
		ORDER BY id
	`
	rows, err := db.Query(getUsersSQL)
	require.NoError(t, rows.Err())
	require.Nil(t, err)

	var actual []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(&u.ID, &u.Login, &u.Pwd)
		require.Nil(t, err)
		actual = append(actual, u)
	}

	require.Equal(t, expected, actual)
}
