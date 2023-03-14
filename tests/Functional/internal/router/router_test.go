package router

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/routing"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPingRoute(t *testing.T) {
	//todo replace with database from env var
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:54323/postgres")
	require.Nil(t, err)
	defer db.Close()

	dbStorage, err := storage.NewDBStorage(db)
	require.Nil(t, err)
	logger := zerolog.Nop()
	r := routing.SetupRouter(dbStorage, &logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Equal(t, `{}`, w.Body.String())
}

func TestRegisterRoute(t *testing.T) {
	db := utils.GetDB(t)
	defer db.Close()
	utils.TruncateTable(t, db)

	dbStorage, err := storage.NewDBStorage(db)
	require.Nil(t, err)
	logger := zerolog.Nop()
	r := routing.SetupRouter(dbStorage, &logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/register", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	assertUser(t, db, "cheesecake")
}

func assertUser(t *testing.T, db *sql.DB, login string) {
	getOneSQL := `
		SELECT id, login, pwd
		FROM "user"
		WHERE login = $1
	`
	row := db.QueryRow(getOneSQL, login)
	require.Nil(t, row.Err())

	var idFromDb int
	var loginFromDB, pwdFromDB string
	err := row.Scan(&idFromDb, &loginFromDB, &pwdFromDB)
	require.Nil(t, err)
}
