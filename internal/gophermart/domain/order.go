package domain

import (
	"context"
	"fmt"
	"github.com/dnsoftware/gophermart/internal/constants"
	"github.com/dnsoftware/gophermart/internal/logger"
	"github.com/dnsoftware/gophermart/internal/storage"
	"strconv"
	"time"
)

type OrderStorage interface {
	Create(ctx context.Context, userID, number int64) (int64, int, error)
	List(ctx context.Context, userID int64) ([]storage.OrderRow, int, error)
}

type Order struct {
	storage OrderStorage
}

// OrderItem plain structure
type OrderItem struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func NewOrderModel(storage OrderStorage) (*Order, error) {
	order := &Order{
		storage: storage,
	}

	return order, nil
}

func (o *Order) AddOrder(ctx context.Context, userID, number int64) (int, error) {
	status := constants.OrderInternalError

	// проверка Луна
	if !IsLuhnValid(number) {
		status = constants.OrderBadNumberFormat
		return status, fmt.Errorf("неверный номер заказа")
	}

	// сохраняем в базу
	id, status, err := o.storage.Create(ctx, userID, number)
	if err != nil {
		return status, err
	}

	logger.Log().Info(fmt.Sprintf("Добавлен заказ %v", id))

	return status, nil
}

func (o *Order) OrdersList(ctx context.Context, userID int64) ([]OrderItem, int, error) {
	orders := make([]OrderItem, 0)

	// получаем список заказов
	list, status, err := o.storage.List(ctx, userID)
	if err != nil {
		return []OrderItem{}, status, err
	}

	for _, item := range list {
		upAt := item.UploadedAt.Format(time.RFC3339)
		orderItem := OrderItem{
			Number:     strconv.FormatInt(item.Num, 10),
			Status:     item.Status,
			Accrual:    item.Accrual,
			UploadedAt: upAt,
		}

		orders = append(orders, orderItem)
	}

	return orders, status, nil
}
