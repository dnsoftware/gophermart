package handlers

import (
	"context"
	"github.com/dnsoftware/gophermart/internal/constants"
	"github.com/dnsoftware/gophermart/internal/gophermart/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UserMart interface {
	AddUser(ctx context.Context, login string, password string) (string, int, error)
	LoginUser(ctx context.Context, login string, password string) (string, int, error)
}

type OrderMart interface {
	AddOrder(ctx context.Context, userID int64, orderID int64) (int, error)
	OrdersList(ctx context.Context, userID int64) ([]domain.OrderItem, int, error)
}

type BalanceMart interface {
	AddTransaction(ctx context.Context, orderNumber int64, amount float32) error
}

type Server struct {
	userMart    UserMart
	orderMart   OrderMart
	balanceMart BalanceMart
	Router      chi.Router
}

type (
	// структура для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// расширенный ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func NewServer(runAddr string, userMart UserMart, orderMart OrderMart, balanceMart BalanceMart) *http.Server {
	h := Server{
		userMart:    userMart,
		orderMart:   orderMart,
		balanceMart: balanceMart,
		Router:      NewRouter(),
	}
	h.Router.Use(trimEnd)
	h.Router.Use(GzipMiddleware)
	h.Router.Use(WithLogging)

	h.Router.Post(constants.UserRegisterAction, h.userRegister)
	h.Router.Post(constants.UserLoginAction, h.userLogin)
	h.Router.With(AuthMiddleware).Post(constants.UserOrderUpload, h.userOrderUpload)
	h.Router.With(AuthMiddleware).Get(constants.UserOrdersList, h.userOrdersList)
	//h.Router.With(AuthMiddleware).Get(constants.UserBalance, h.userBalance)

	srv := &http.Server{
		Addr:    runAddr,
		Handler: h.Router,
	}

	return srv
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
