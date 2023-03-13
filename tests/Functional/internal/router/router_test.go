package router

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	sut "github.com/smamykin/gofermart/internal/router"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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
	router := sut.SetupRouter(dbStorage, &logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Equal(t, `{}`, w.Body.String())
}

//func TestRegisterRoute(t *testing.T) {
//	router := sut.SetupRouter()
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/api/user/register", strings.NewReader(`{"login":"cheesecake", "password": "pancake"}`))
//	router.ServeHTTP(w, req)
//
//	require.Equal(t, 200, w.Code)
//	require.Equal(t, "pong", w.Body.String())
//}
