package router

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/routing"
	"github.com/smamykin/gofermart/pkg/pwdhash"
	"github.com/smamykin/gofermart/pkg/token"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	db := utils.GetDB(t)
	defer db.Close()

	logger := zerolog.Nop()
	r := routing.SetupRouter(db, &logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, `{"DBError":""}`, w.Body.String())
	require.Equal(t, 200, w.Code)
}

func TestRegister(t *testing.T) {
	db := utils.GetDB(t)
	defer db.Close()
	utils.TruncateTable(t, db)

	logger := zerolog.Nop()
	r := routing.SetupRouter(db, &logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/register", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	assertUser(t, db, "cheesecake")
}

func TestLogin(t *testing.T) {
	db := utils.GetDB(t)
	defer db.Close()
	utils.TruncateTable(t, db)

	pwd := "pancake"
	login := "cheesecake"
	addUserToDB(t, pwd, login, db)

	logger := zerolog.Nop()
	r := routing.SetupRouter(db, &logger)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/login", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, w.Body.String())
	require.Equal(t, `{"message":"success"}`, w.Body.String())

	bearerToken := w.Header().Get("Authorization")
	require.NotSame(t, "", bearerToken)
	require.Equal(t, 2, len(strings.Split(bearerToken, " ")))

	tokenString := strings.Split(bearerToken, " ")[1]
	tkn, err := token.ParseString(tokenString)
	require.Nil(t, err)
	require.Equal(t, true, tkn.Valid)

	claims, _ := tkn.Claims.(jwt.MapClaims)
	id, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	require.Equal(t, 1, int(id))

}

func addUserToDB(t *testing.T, pwd string, login string, db *sql.DB) {
	hg := pwdhash.HashGenerator{}
	pwdHash, err := hg.Generate(pwd)
	require.Nil(t, err)
	_, err = db.Exec(`INSERT INTO "user" (login, pwd) VALUES ($1, $2)`, login, pwdHash)
	require.Nil(t, err)
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
