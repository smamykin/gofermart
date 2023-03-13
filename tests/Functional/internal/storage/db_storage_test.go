package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDBStorage_UpsertUser(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:54323/postgres")
	require.Nil(t, err)
	defer db.Close()

	truncateTable(t, db)
	store, err := storage.NewDBStorage(db)
	require.Nil(t, err)

	err = store.UpsertUser("cheesecake", "pwd")
	require.Nil(t, err)
	assertUsersInDB(t, db, []entity.User{
		{1, "cheesecake", "pwd"},
	})

	err = store.UpsertUser("cheesecake", "pwd2")
	assertUsersInDB(t, db, []entity.User{
		{1, "cheesecake", "pwd2"},
	})
	err = store.UpsertUser("cheesecake2", "pwd")
	assertUsersInDB(t, db, []entity.User{
		{1, "cheesecake", "pwd2"},
		{3, "cheesecake2", "pwd"},
	})
}

func TestDBStorage_GetUserByLogin(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:54323/postgres")
	require.Nil(t, err)
	defer db.Close()

	truncateTable(t, db)
	expected := entity.User{
		ID:    1,
		Login: "foo",
		Pwd:   "bar",
	}
	insertUser(t, db, expected)

	store, err := storage.NewDBStorage(db)
	require.Nil(t, err)
	actual, err := store.GetUserByLogin("foo")
	require.Nil(t, err)
	assert.Equal(t, expected, actual)
	_, err = store.GetUserByLogin("baz")
	assert.Equal(t, service.ErrNoRows, err)

}

func insertUser(t *testing.T, db *sql.DB, expected entity.User) {
	_, err := db.Exec(`INSERT INTO "user" (login, pwd) VALUES ($1,$2);`, expected.Login, expected.Pwd)
	require.Nil(t, err)
}

func truncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`TRUNCATE TABLE public."user" RESTART IDENTITY;`)
	require.Nil(t, err)
}

func assertUsersInDB(t *testing.T, db *sql.DB, expected []entity.User) {
	getUsersSQL := `
		SELECT id, login, pwd
		FROM public."user"
		ORDER BY id
	`
	rows, err := db.Query(getUsersSQL)
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
