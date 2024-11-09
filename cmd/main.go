package main

import (
	"fmt"
	"github.com/Davmie/javaCode/cmd/server"
	walletDel "github.com/Davmie/javaCode/internal/wallet/delivery"
	pgWallet "github.com/Davmie/javaCode/internal/wallet/repository/postgres"
	walletUseCase "github.com/Davmie/javaCode/internal/wallet/usecase"
	"github.com/Davmie/javaCode/pkg/middleware"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=postgres user=program password=test dbname=wallets port=5432"}

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	walletHandler := walletDel.WalletHandler{
		WalletUseCase: walletUseCase.New(pgWallet.New(logger, db)),
		Logger:        logger,
	}

	r := http.NewServeMux()

	r.Handle("POST /api/v1/wallets", http.HandlerFunc(walletHandler.Create))
	//r.Handle("GET /api/v1/wallets/{walletId}", http.HandlerFunc(walletHandler.Get))
	r.Handle("PATCH /api/v1/wallets/{walletId}", http.HandlerFunc(walletHandler.Update))
	r.Handle("DELETE /api/v1/wallets/{walletId}", http.HandlerFunc(walletHandler.Delete))
	r.Handle("GET /api/v1/wallets", http.HandlerFunc(walletHandler.GetAll))
	r.Handle("POST /api/v1/wallet", http.HandlerFunc(walletHandler.ChangeAmount))
	r.Handle("GET /api/v1/wallets/{WALLET_UUID}", http.HandlerFunc(walletHandler.GetByUID))

	router := middleware.AccessLog(logger, r)
	router = middleware.Panic(logger, router)

	s := server.NewServer(router)
	if err := s.Start(); err != nil {
		logger.Fatal(err)
	}

	err = zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
