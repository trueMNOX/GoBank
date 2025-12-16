package http

import (
	"Gobank/internal/service"
	"Gobank/internal/token"
	"Gobank/internal/transport/http/handler"
	"Gobank/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	tokenMaker token.TokenMaker,
	authService service.AuthService,
	accountService service.AccountService,
	transferService service.TransferService,
) *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	transferHandler := handler.NewTransferHandler(transferService)

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	authRoutes := api.Group("/").Use(middleware.AuthMiddleware(tokenMaker))
	{
		//account
		authRoutes.POST("/accounts", accountHandler.CreateAccount)
		authRoutes.GET("/accounts/:id", accountHandler.GetAccount)
		authRoutes.GET("/accounts", accountHandler.ListAccount)

		//transfer
		authRoutes.POST("/transfer", transferHandler.CreateTransfer)
	}
	return r
}
