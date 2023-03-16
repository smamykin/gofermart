package utils

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/container"
	"github.com/stretchr/testify/require"
	"testing"
)

func TruncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`TRUNCATE TABLE "user" RESTART IDENTITY;`)
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
