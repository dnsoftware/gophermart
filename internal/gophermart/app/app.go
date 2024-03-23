package app

import (
	"context"
	"fmt"
	"github.com/dnsoftware/gophermart/internal/gophermart/config"
	"github.com/dnsoftware/gophermart/internal/gophermart/domain"
	"github.com/dnsoftware/gophermart/internal/gophermart/handlers"
	"github.com/dnsoftware/gophermart/internal/logger"
	"github.com/dnsoftware/gophermart/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run() {
	var wg sync.WaitGroup
	gLogger := logger.Log()
	defer gLogger.Sync()

	cfg := config.NewServerConfig()

	martStorage, err := storage.NewMartStorage(cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}

	userRepo := storage.NewUserRepo(martStorage)
	orderRepo := storage.NewOrderRepo(martStorage)
	balanceRepo := storage.NewBalanceRepo(martStorage)
	accrualRepo := storage.NewAccrualRepo()

	gophermart, err := domain.NewGophermart(cfg, userRepo, orderRepo, balanceRepo, accrualRepo)
	if err != nil {
		panic(err)
	}

	ctxSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// канал с ордерами на проверку
	chanUnchecked := domain.NewOrdersUnchecked()

	// канал с проверенными ордерами
	chanChecked := domain.NewOrdersChecked()

	// работа с accrual сервисом - проверка начислений по заказам
	wg.Add(1)
	go func() {
		defer wg.Done()
		gophermart.Accrual.StartAccrualChecker(ctxSignal)
	}()

	srv := handlers.NewServer(cfg.RunAddress, gophermart.User, gophermart.Order, gophermart.Balance)

	// запуск HTTP сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Listening and serving")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// ждет сигнала завершения программы, потом завершает работу HTTP сервера и работу с Accrual
	wg.Add(1)
	go func() {
		<-ctxSignal.Done()

		fmt.Println("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			wg.Done()
			stop()
			cancel()
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			logger.Log().Error("Ошибка при завершении HTTP сервера: " + err.Error())
			return
		}

		fmt.Println("Shutdown completed")
	}()

	wg.Wait()

	fmt.Println("\nПрограмма завершена!")

}
