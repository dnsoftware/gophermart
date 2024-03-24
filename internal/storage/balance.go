package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dnsoftware/gophermart/internal/constants"
	"time"
)

type BalanceRepo struct {
	storage *MartStorage
}

type BalanceRow struct {
	ID          int64
	UserID      int64
	OrderNumber int64
	Amount      float32
	ProcessedAt time.Time
}

func NewBalanceRepo(storage *MartStorage) *BalanceRepo {

	balance := BalanceRepo{
		storage: storage,
	}

	return &balance
}

func (b *BalanceRepo) SaveTransaction(ctx context.Context, orderNumber int64, amount float32) error {

	// получение ID владельца заказа
	var userID int64
	query := `SELECT user_id 
			  FROM orders WHERE num = $1`
	row := b.storage.db.QueryRowContext(ctx, query, orderNumber)

	err := row.Scan(&userID)
	if err != nil {
		return err
	}

	// проверка, что данные по этому заказу еще не вносились в базу
	var a float32
	query = `SELECT amount 
			 FROM balances WHERE order_number = $1`
	row = b.storage.db.QueryRowContext(ctx, query, orderNumber)

	err = row.Scan(&a)
	if err != nil && err != sql.ErrNoRows { // sql.ErrNoRows если нет строки
		return err
	}
	if err == nil { // если ошибок нет - то такая запись уже есть
		return fmt.Errorf("Данные по начислению уже сохранены")
	}

	// стартуем транзакцию БД
	tx, err := b.storage.db.Begin()
	if err != nil {
		return err
	}

	// обновление статуса ордера на обработанный
	query = `UPDATE orders SET status = $1, accrual = $2
			  WHERE num = $3`

	err = b.storage.retryExec(ctx, query, constants.OrderProcessed, amount, orderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	// занесение начислений на баланс
	query = `INSERT INTO balances (user_id, order_number, amount, processed_at)
			  VALUES ($1, $2, $3, now())`
	err = b.storage.retryExec(ctx, query, userID, orderNumber, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
