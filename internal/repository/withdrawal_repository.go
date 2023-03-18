package repository

import (
	"database/sql"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"time"
)

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

type WithdrawalRepository struct {
	db *sql.DB
}

func (w *WithdrawalRepository) GetWithdrawal(ID int) (order entity.Withdrawal, err error) {
	row := w.db.QueryRow(`
		SELECT id, user_id, amount, created_at
		FROM "withdrawal"
		WHERE id = $1
	`, ID)

	return hydrateWithdrawal(row)
}

func (w *WithdrawalRepository) GetAmountSumByUserId(userID int) (sum float64, err error) {
	row := w.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0.00)
		FROM withdrawal
		WHERE user_id = $1
	`, userID)

	if row.Err() != nil {
		return 0, row.Err()
	}

	err = row.Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
func (w *WithdrawalRepository) AddWithdrawal(withdrawal entity.Withdrawal) (entity.Withdrawal, error) {
	withdrawal.CreatedAt = time.Now().UTC().Truncate(time.Second)
	row := w.db.QueryRow(
		`
			INSERT INTO "withdrawal" (user_id, amount, created_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
		withdrawal.UserID,
		withdrawal.Amount,
		withdrawal.CreatedAt,
	)

	if row.Err() != nil {
		return withdrawal, row.Err()
	}

	err := row.Scan(&withdrawal.ID)

	return withdrawal, err
}

func hydrateWithdrawal(row *sql.Row) (withdrawal entity.Withdrawal, err error) {
	if row.Err() != nil {
		return withdrawal, row.Err()
	}

	err = row.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.Amount, &withdrawal.CreatedAt)

	if err == sql.ErrNoRows {
		return withdrawal, service.ErrEntityIsNotFound
	}

	return withdrawal, err
}
