package domain

import "github.com/dnsoftware/gophermart/internal/constants"

/* Очередь проверенных ордеров на занесение в базу */

type orderData struct {
	order   int64
	status  string
	accrual float32
}

type OrdersChecked struct {
	ordersCh chan orderData
}

func NewOrdersChecked() *OrdersChecked {
	return &OrdersChecked{
		ordersCh: make(chan orderData, constants.OrdersChannelCapacity),
	}
}

// ставим в очередь для дальнейшего сохранения в базу
func (c *OrdersChecked) Push(order int64, status string, accrual float32) {

	o := orderData{
		order:   order,
		status:  status,
		accrual: accrual,
	}

	c.ordersCh <- o
}

// забираем из очереди для сохранения в базу
func (c *OrdersChecked) Pop() (int64, string, float32) {

	o := <-c.ordersCh

	return o.order, o.status, o.accrual
}
