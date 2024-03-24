package domain

import (
	"context"
	"fmt"
)

type BalanceStorage interface {
	SaveTransaction(ctx context.Context, orderNumber int64, amount float32) error
}

type Balance struct {
	storage BalanceStorage
}

func NewBalanceModel(storage BalanceStorage) *Balance {
	balance := &Balance{
		storage: storage,
	}

	return balance
}

// комплексное обновление данных в базе
func (b *Balance) AddTransaction(ctx context.Context, orderNumber int64, amount float32) error {

	err := b.storage.SaveTransaction(ctx, orderNumber, amount)
	if err != nil {
		return fmt.Errorf("Ошибка при сохранении начислений в базу %w", err)
	}

	return nil
}
