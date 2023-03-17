package repository

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

func TestUserRepository_UpsertUser(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)
	repository := c.UserRepository()

	// insert
	user, err := repository.UpsertUser("cheesecake", "pwd")
	require.NoError(t, err)
	expectedUser := entity.User{ID: 1, Login: "cheesecake", Pwd: "pwd"}
	assertUsersInDB(t, db, []entity.User{expectedUser})
	require.Equal(t, expectedUser, user)

	//update the same
	user, err = repository.UpsertUser("cheesecake", "pwd2")
	require.NoError(t, err)
	expectedUser = entity.User{ID: 1, Login: "cheesecake", Pwd: "pwd2"}
	assertUsersInDB(t, db, []entity.User{expectedUser})
	require.Equal(t, expectedUser, user)

	//add a new user again
	user, err = repository.UpsertUser("cheesecake2", "pwd")
	require.NoError(t, err)
	expectedUser = entity.User{ID: 3, Login: "cheesecake2", Pwd: "pwd"}
	assertUsersInDB(t, db, []entity.User{
		{ID: 1, Login: "cheesecake", Pwd: "pwd2"},
		{ID: 3, Login: "cheesecake2", Pwd: "pwd"},
	})
	require.Equal(t, expectedUser, user)
}

func TestUserRepository_GetUserByLogin(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	expected := utils.InsertUser(t, db, entity.User{
		Login: "foo",
		Pwd:   "bar",
	})

	repository := c.UserRepository()

	actual, err := repository.GetUserByLogin("foo")
	require.Nil(t, err)
	assert.Equal(t, expected, actual)

	_, err = repository.GetUserByLogin("baz")
	assert.Equal(t, service.ErrEntityIsNotFound, err)

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
