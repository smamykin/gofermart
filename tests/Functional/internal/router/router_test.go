package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/smamykin/gofermart/internal/container"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/pkg/money"
	"github.com/smamykin/gofermart/pkg/pwdhash"
	"github.com/smamykin/gofermart/pkg/token"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
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
	addUserWithHashedPwdToDB(t, pwd, login, c.DB())

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
	userID := addUserWithHashedPwdToDB(t, pwd, login, c.DB()).ID

	r := c.Router()
	w := httptest.NewRecorder()
	orderNumber := "12345678903"
	req, _ := http.NewRequest("POST", "/api/user/orders", strings.NewReader(orderNumber))
	authorize(t, userID, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusAccepted, w.Code, w.Body.String())
	actualOrder := entity.Order{}
	err := json.Unmarshal(w.Body.Bytes(), &actualOrder)
	require.NoError(t, err)
	require.Equal(t, actualOrder.OrderNumber, orderNumber)
	require.Equal(t, actualOrder.UserID, userID)
}

func TestOrderList(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	user := utils.InsertUser(t, c.DB(), entity.User{})
	userNotToGet := utils.InsertUser(t, c.DB(), entity.User{})

	//create orders to get
	orderToGet, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "123",
	})
	require.NoError(t, err)
	_, err = c.OrderRepository().AddOrder(entity.Order{
		UserID:      userNotToGet.ID,
		OrderNumber: "321",
	})
	require.NoError(t, err)

	r := c.Router()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/orders", nil)
	authorize(t, user.ID, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, w.Body.String())
	var actualResponse []controller.OrderResponseModel
	err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)
	require.Equal(t, []controller.OrderResponseModel{
		{
			Number:     orderToGet.OrderNumber,
			Accrual:    orderToGet.Accrual.AsFloat(),
			Status:     orderToGet.Status.String(),
			UploadedAt: orderToGet.CreatedAt.Format(time.RFC3339),
		},
	}, actualResponse)
}

func TestBalance(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	user := utils.InsertUser(t, c.DB(), entity.User{})
	_, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:        user.ID,
		OrderNumber:   "123",
		Status:        entity.OrderStatusProcessed,
		AccrualStatus: entity.AccrualStatusProcessed,
		Accrual:       money.FromFloat(442.5),
	})
	require.NoError(t, err)
	_, err = c.WithdrawalRepository().AddWithdrawal(entity.Withdrawal{
		UserID: user.ID,
		Amount: money.FromFloat(42.3),
	})
	require.NoError(t, err)

	r := c.Router()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/balance", nil)
	authorize(t, user.ID, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, `{"current":400.2,"withdrawn":42.3}`, w.Body.String())
}

func TestWithdraw(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	orderNumber := "12345678903"
	sum := 40.4

	user := utils.InsertUser(t, c.DB(), entity.User{})
	_, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "123",
		Accrual:     money.FromFloat(sum),
	})
	require.NoError(t, err)

	r := c.Router()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/balance/withdraw", strings.NewReader(fmt.Sprintf(`{"order":"%s","sum":%f}`, orderNumber, sum)))
	authorize(t, user.ID, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	actualWithdrawal := entity.Withdrawal{}
	err = json.Unmarshal(w.Body.Bytes(), &actualWithdrawal)
	require.NoError(t, err)
	require.Equal(t, orderNumber, actualWithdrawal.OrderNumber)
	require.Equal(t, user.ID, actualWithdrawal.UserID)
}

func TestWithdrawalList(t *testing.T) {
	c := utils.GetContainer(t)
	utils.TruncateTable(t, c.DB())

	user := utils.InsertUser(t, c.DB(), entity.User{})
	userNotToGet := utils.InsertUser(t, c.DB(), entity.User{})

	//create orders to get
	withdrawalToGet, err := c.WithdrawalRepository().AddWithdrawal(entity.Withdrawal{
		UserID:      user.ID,
		OrderNumber: "123",
		Amount:      111,
	})
	require.NoError(t, err)
	_, err = c.WithdrawalRepository().AddWithdrawal(entity.Withdrawal{
		UserID:      userNotToGet.ID,
		OrderNumber: "321",
	})
	require.NoError(t, err)

	r := c.Router()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/withdrawals", nil)
	authorize(t, user.ID, c, req)

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, w.Body.String())
	var actualResponse []controller.WithdrawalResponseModel
	err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)
	require.Equal(t, []controller.WithdrawalResponseModel{
		{
			OrderNumber: withdrawalToGet.OrderNumber,
			Amount:      withdrawalToGet.Amount.AsFloat(),
			ProcessedAt: withdrawalToGet.CreatedAt.Format(time.RFC3339),
		},
	}, actualResponse)
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

func addUserWithHashedPwdToDB(t *testing.T, pwd string, login string, db *sql.DB) entity.User {
	hg := pwdhash.HashGenerator{}
	pwdHash, err := hg.Generate(pwd)
	require.Nil(t, err)

	return utils.InsertUser(t, db, entity.User{
		Login: login, Pwd: pwdHash,
	})
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
