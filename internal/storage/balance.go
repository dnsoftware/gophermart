package storage

import (
	"context"
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

func (b *BalanceRepo) Create(ctx context.Context, userID, number int64) (int64, int, error) {

	return 0, 0, nil
}
