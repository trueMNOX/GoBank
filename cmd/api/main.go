package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Gobank/internal/database"
	"Gobank/internal/repository"
	"Gobank/internal/service"
	"Gobank/internal/token"
	httptransport "Gobank/internal/transport/http"
	"Gobank/pkg/config"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	defer func() {
		database.CloseDatabase(db)
	}()

	tokenMaker, err := token.NewJwtMaker(cfg.JWTSecret)
	if err != nil {
		log.Fatal("cannot create token maker: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	entryRepo := repository.NewEntryRepositoryImpl(db)
	transferRepo := repository.NewTransferRepository(db, entryRepo, accountRepo)

	authService, err := service.NewAuthService(userRepo, cfg)
	if err != nil {
		log.Fatal("cannot create auth service: %v", err)
	}
	accountService := service.NewAccountService(accountRepo, userRepo)
	transferService := service.NewTransferService(transferRepo, db, accountRepo)

	r := httptransport.SetupRouter(tokenMaker, *authService, *accountService, *transferService)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Println("starting server on port: %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server exiting")
	}
}
