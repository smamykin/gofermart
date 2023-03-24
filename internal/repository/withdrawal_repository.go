package repository

import (
	"database/sql"
	"errors"
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

func (w *WithdrawalRepository) GetAllByUserID(userID int) ([]entity.Withdrawal, error) {
	rows, err := w.db.Query(
		`
			SELECT id, user_id, order_number, amount, created_at
			FROM "withdrawal"
			WHERE user_id = $1
			ORDER BY created_at
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return hydrateWithdrawals(rows)
}

func (w *WithdrawalRepository) GetWithdrawalByOrderNumber(orderNumber string) (order entity.Withdrawal, err error) {
	row := w.db.QueryRow(
		`
			SELECT id, user_id, order_number, amount, created_at
			FROM "withdrawal"
			WHERE order_number = $1
		`,
		orderNumber,
	)

	return hydrateWithdrawal(row)
}

func (w *WithdrawalRepository) GetWithdrawal(ID int) (order entity.Withdrawal, err error) {
	row := w.db.QueryRow(`
		SELECT id, user_id, order_number, amount, created_at
		FROM "withdrawal"
		WHERE id = $1
	`, ID)

	return hydrateWithdrawal(row)
}

func (w *WithdrawalRepository) GetAmountSumByUserID(userID int) (sum float64, err error) {
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
			INSERT INTO "withdrawal" (user_id, order_number, amount,created_at)
			VALUES ($1,$2,$3,$4)
			RETURNING id
		`,
		withdrawal.UserID,
		withdrawal.OrderNumber,
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

	err = row.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderNumber, &withdrawal.Amount, &withdrawal.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return withdrawal, service.ErrEntityIsNotFound
	}

	return withdrawal, err
}

func hydrateWithdrawals(rows *sql.Rows) (withdrawals []entity.Withdrawal, err error) {
	withdrawals = make([]entity.Withdrawal, 0)
	for rows.Next() {
		var withdrawal entity.Withdrawal
		err = rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderNumber, &withdrawal.Amount, &withdrawal.CreatedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
