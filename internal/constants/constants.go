package constants

import (
	"go.uber.org/zap/zapcore"
	"time"
)

// параметры работы сервера по умолчанию
const (
	RunAddress      string = "localhost:8081" // адрес:порт сервера по умолчанию
	AccrualAddress  string = "localhost:8080" // адрес:порт сервера по умолчанию
	AccrualProtocol string = "http"
)

// интервалы
const (
	DBContextTimeout   time.Duration = time.Duration(50) * time.Second // длительность запроса в контексте работы с БД
	HTTPContextTimeout time.Duration = time.Duration(10) * time.Second // длительность запроса в контексте работы с сетью
	HTTPAttemtPeriods  string        = "1s,2s,5s"
)

// логгер
const (
	LogFile  string = "log.log"
	LogLevel        = zapcore.InfoLevel
)

const (
	EncodingGzip   string = "gzip"
	HashHeaderName string = "HashSHA256"
)

// crypto
const (
	PasswordSalt string = "gT65HtksdhHHj"
)

// routes
const (
	UserRegisterRoute    string = "/api/user/register"
	UserLoginRoute       string = "/api/user/login"
	UserOrderUploadRoute string = "/api/user/orders"
	UserOrdersListRoute  string = "/api/user/orders"
	UserBalanceRoute     string = "/api/user/balance"
	UserWithdrawalsRoute string = "/api/user/withdrawals"
	UserWithdrawRoute    string = "/api/user/balance/withdraw"
)

// разное
const (
	ApplicationJSON string = "application/json"

	TokenExp  = time.Hour * 24
	SecretKey = "golangforever"

	MinLoginLength    = 3
	MinPasswordLength = 8

	HeaderAuthorization = "Authorization"
	UserIDKey           = "userID"

	AccrualServiceQueryLimit = 100                    // максимально кол-во запросов к Accrual сервису в минуту
	AccrualCheckPeriod       = 1                      // период проверки
	AccrualOrderEndpoint     = "/api/orders/{number}" // получение информации о расчёте начислений баллов лояльности

	OrdersChannelCapacity = 10 // емкость канала для обмена данными по ордерам
	CheckOrdersPeriod     = 1  // период проверки необработанных ордеров в секундах
)

// статусы заказов
const (
	OrderNew        = "NEW"
	OrderProcessing = "PROCESSING"
	OrderInvalid    = "INVALID"
	OrderProcessed  = "PROCESSED"
)

// статусы расчетов
const (
	AccrualRegistered = "REGISTERED"
	AccrualInvalid    = "INVALID"
	AccrualProcessing = "PROCESSING"
	AccrualProcessed  = "PROCESSED"
)
