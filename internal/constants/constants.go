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
	DBContextTimeout   time.Duration = time.Duration(5) * time.Second  // длительность запроса в контексте работы с БД
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

// actions
const (
	UserRegisterAction string = "/api/user/register"
	UserLoginAction    string = "/api/user/login"
	UserOrderUpload    string = "/api/user/orders"
	UserOrdersList     string = "/api/user/orders"
)

// разное
const (
	ApplicationJSON string = "application/json"

	TOKEN_EXP  = time.Hour * 24
	SECRET_KEY = "golangforever"

	MinLoginLength    = 3
	MinPasswordLength = 8

	HeaderAuthorization = "Authorization"
	UserIDKey           = "userID"

	AccrualServiceQueryLimit = 100                    // максимально кол-во запросов к Accrual сервису в минуту
	AccrualCheckPeriod       = 10                     // период проверки
	AccrualOrderEndpoint     = "/api/orders/{number}" // получение информации о расчёте начислений баллов лояльности

	OrdersChannelCapacity = 1 // емкость канала для обмена данными по ордерам
	CheckOrdersPeriod     = 5 // период проверки необработанных ордеров в секундах
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