package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dnsoftware/gophermart/internal/constants"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type AccrualRepo struct {
	client                *http.Client
	orderEndpointTemplate string
}

type AccrualRow struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func NewAccrualRepo() *AccrualRepo {
	return &AccrualRepo{
		client:                &http.Client{},
		orderEndpointTemplate: constants.AccrualProtocol + "://" + constants.AccrualAddress + constants.AccrualOrderEndpoint,
	}
}

// Получить данные по заказу
func (a *AccrualRepo) GetOrder(orderNum int64) (*AccrualRow, int, error) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	url := a.orderEndpoint(orderNum)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, buf)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	request.Header.Set("Content-Type", constants.ApplicationJSON)

	resp, err := a.client.Do(request)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	row := &AccrualRow{}
	err = json.Unmarshal(data, row)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return row, resp.StatusCode, nil
}

func (a *AccrualRepo) orderEndpoint(orderNum int64) string {
	strNum := strconv.FormatInt(orderNum, 10)

	return strings.Replace(a.orderEndpointTemplate, "{number}", strNum, -1)
}
