package domain

import (
	"context"
	"fmt"
	"github.com/dnsoftware/gophermart/internal/constants"
	"github.com/dnsoftware/gophermart/internal/logger"
	"github.com/dnsoftware/gophermart/internal/storage"
	"time"
)

type AccrualStorage interface {
	GetOrder(orderNum int64) (*storage.AccrualRow, int, error)
}

type AccrualItem struct {
}

type Accrual struct {
	storage                  AccrualStorage
	accrualServiceQueryLimit int           // максимально кол-во запросов к Accrual сервису в минуту
	counter                  int           // счетчик запросов за текущий интервал
	checkPeriod              time.Duration // период проверки
}

func NewAccrualModel(storage AccrualStorage) (*Accrual, error) {
	balance := &Accrual{
		storage:                  storage,
		accrualServiceQueryLimit: constants.AccrualServiceQueryLimit,
		counter:                  0,
		checkPeriod:              time.Duration(constants.AccrualCheckPeriod) * time.Second,
	}

	return balance, nil
}

// StartAccrualChecker Служба проверки начислений
func (b *Accrual) StartAccrualChecker(ctx context.Context) {
	timer := time.NewTimer(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			logger.Log().Info("AccrualChecker DONE!!!")
			return
		case <-timer.C:
			b.counter = 0
			fmt.Println("tick")
		default:
			if b.counter == b.accrualServiceQueryLimit {
				timer.Reset(b.checkPeriod)
				continue
			}

			// основная работа
			fmt.Println("accepting request", b.counter)
			// TODO запрос к каналу ордеров и получение необработанных
			b.accrualCheck(2224764148437)

			b.counter++
		}
	}

}

func (b *Accrual) accrualCheck(orderNum int64) {
	order, status, err := b.storage.GetOrder(orderNum)
	if err != nil {

	}

	fmt.Println(order, status)
}
