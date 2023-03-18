package repository

import (
	"database/sql"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"time"
)

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderRepository struct {
	db *sql.DB
}

func (o *OrderRepository) AddOrder(order entity.Order) (entity.Order, error) {

	order.CreatedAt = time.Now().UTC().Truncate(time.Second)
	row := o.db.QueryRow(
		`
			INSERT INTO "order" (user_id, order_number, status, accrual_status, accrual, created_at) 
			VALUES ($1, $2, $3, $4,$5, $6)
			RETURNING id
		`,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.AccrualStatus,
		order.Accrual,
		order.CreatedAt,
	)

	if row.Err() != nil {
		return order, row.Err()
	}

	err := row.Scan(&order.ID)

	return order, err
}

func (o *OrderRepository) GetOrder(ID int) (order entity.Order, err error) {
	row := o.db.QueryRow(`
		SELECT id, user_id, order_number, status, accrual_status, accrual, created_at
		FROM "order"
		WHERE id = $1
	`, ID)

	return hydrateOrder(row)
}

func (o *OrderRepository) GetOrderByOrderNumber(orderNumber string) (order entity.Order, err error) {
	row := o.db.QueryRow(
		`
			SELECT id, user_id, order_number, status, accrual_status, accrual, created_at
			FROM "order"
			WHERE order_number = $1
		`,
		orderNumber,
	)

	return hydrateOrder(row)
}

func (o *OrderRepository) GetAllByUserID(userID int) ([]entity.Order, error) {
	rows, err := o.db.Query(
		`
			SELECT id, user_id, order_number, status, accrual_status, accrual, created_at
			FROM "order"
			WHERE user_id = $1
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return hydrateOrders(rows)
}

func (o *OrderRepository) UpdateOrder(order entity.Order) (entity.Order, error) {
	order.CreatedAt = time.Now().UTC().Truncate(time.Second)
	result, err := o.db.Exec(
		`
			UPDATE "order"
			SET status = $1,
			    accrual_status = $2,
			    accrual = $3,
				user_id = $4
			WHERE order_number = $5
		`,
		//set
		order.Status,
		order.AccrualStatus,
		order.Accrual,
		order.UserID,
		//where
		order.OrderNumber,
	)

	if err != nil {
		return order, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return order, err
	}

	if affected == 0 {
		return order, service.ErrEntityIsNotFound
	}

	return order, nil
}

func (o *OrderRepository) GetOrdersWithUnfinishedStatus() ([]entity.Order, error) {
	rows, err := o.db.Query(
		`
			SELECT id, user_id, order_number, status, accrual_status, accrual, created_at
			FROM "order"
			WHERE status in ($1, $2)
		`,
		entity.OrderStatusNew,
		entity.OrderStatusProcessing,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return hydrateOrders(rows)
}

func hydrateOrder(row *sql.Row) (order entity.Order, err error) {
	if row.Err() != nil {
		return order, row.Err()
	}

	err = row.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status, &order.AccrualStatus, &order.Accrual, &order.CreatedAt)

	if err == sql.ErrNoRows {
		return order, service.ErrEntityIsNotFound
	}

	return order, err
}

func hydrateOrders(rows *sql.Rows) (orders []entity.Order, err error) {
	orders = make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		err = rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status, &order.AccrualStatus, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
