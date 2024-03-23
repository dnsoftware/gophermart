package domain

import (
	"context"
)

type BalanceStorage interface {
	Create(ctx context.Context, userID, number int64) (int64, int, error)
}

type Balance struct {
	storage BalanceStorage
}

func NewBalanceModel(storage BalanceStorage) (*Balance, error) {
	balance := &Balance{
		storage: storage,
	}

	return balance, nil
}

func (b *Balance) AddTransaction(ctx context.Context, userID int64, orderNumber int64, amount float32) (int, error) {

	return 0, nil
}
