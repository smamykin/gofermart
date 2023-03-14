package utils

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TruncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`TRUNCATE TABLE "user" RESTART IDENTITY;`)
	require.Nil(t, err)
}

func GetDB(t *testing.T) *sql.DB {
	dsn, ok := os.LookupEnv("DATABASE_DSN")
	require.True(t, ok, "database dsn is not set")
	db, err := sql.Open("pgx", dsn)
	require.Nil(t, err)
	return db
}
