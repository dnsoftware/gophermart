package domain

import "github.com/dnsoftware/gophermart/internal/constants"

/* Очередь ордеров на проверку */

type OrdersUnchecked struct {
	ordersCh chan int64
}

func NewOrdersUnchecked() *OrdersUnchecked {
	return &OrdersUnchecked{
		ordersCh: make(chan int64, constants.OrdersChannelCapacity),
	}
}

func (u *OrdersUnchecked) Push(number int64) {
	u.ordersCh <- number
}

func (u *OrdersUnchecked) Pop() int64 {
	return <-u.ordersCh
}
