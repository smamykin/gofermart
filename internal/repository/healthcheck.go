package repository

import (
	"context"
	"database/sql"
	"time"
)

func CreateHealthcheckFunc(db *sql.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			return err
		}

		return nil
	}
}
