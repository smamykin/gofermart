package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/smamykin/gofermart/internal/container"
	"github.com/smamykin/gofermart/internal/entity"
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
	c := utils.GetContainer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	c.Router().ServeHTTP(w, req)

	require.Equal(t, `{"DBError":""}`, w.Body.String())
	require.Equal(t, 200, w.Code)
}

func TestRegister(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/register", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
	c.Router().ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	assertUser(t, c.DB(), "cheesecake")
	require.Equal(t, `{"message":"success"}`, w.Body.String())
	assertAuthorizationHeader(t, w, c)

}

func TestLogin(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	pwd := "pancake"
	login := "cheesecake"
	addUserToDB(t, pwd, login, c.DB())

	r := c.Router()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/login", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, w.Body.String())
	require.Equal(t, `{"message":"success"}`, w.Body.String())
	assertAuthorizationHeader(t, w, c)
}

func TestOrderPost(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	pwd := "pancake"
	login := "cheesecake"
	addUserToDB(t, pwd, login, c.DB())

	r := c.Router()
	w := httptest.NewRecorder()
	orderNumber := "12345678903"
	req, _ := http.NewRequest("POST", "/api/user/orders", strings.NewReader(orderNumber))
	userId := 1
	authorize(t, userId, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, w.Body.String())
	actualOrder := entity.Order{}
	err := json.Unmarshal(w.Body.Bytes(), &actualOrder)
	require.NoError(t, err)
	require.Equal(t, actualOrder.OrderNumber, orderNumber)
	require.Equal(t, actualOrder.UserID, userId)
}

func authorize(t *testing.T, userID int, c *container.Container, req *http.Request) {
	tkn, err := token.Generate(userID, []byte(c.Config().APISecret), c.Config().TokenLifespan)
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))
}

func assertAuthorizationHeader(t *testing.T, w *httptest.ResponseRecorder, c *container.Container) {
	bearerToken := w.Header().Get("Authorization")
	require.NotSame(t, "", bearerToken)
	require.Equal(t, 2, len(strings.Split(bearerToken, " ")))

	tokenString := strings.Split(bearerToken, " ")[1]
	tkn, err := token.ParseString(tokenString, []byte(c.Config().APISecret))
	require.Nil(t, err)
	require.Equal(t, true, tkn.Valid)

	claims, _ := tkn.Claims.(jwt.MapClaims)
	id, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	require.NoError(t, err)
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

	var idFromDB int
	var loginFromDB, pwdFromDB string
	err := row.Scan(&idFromDB, &loginFromDB, &pwdFromDB)
	require.Nil(t, err)
}
