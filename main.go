package main

import (
	"fmt"

	"github.com/torngkab/grit-api-gateway-service/app/accounts/createaccount"
	"github.com/torngkab/grit-api-gateway-service/app/accounts/getbalance"
	"github.com/torngkab/grit-api-gateway-service/app/accounts/gettransactions"
	"github.com/torngkab/grit-api-gateway-service/app/deposit"
	"github.com/torngkab/grit-api-gateway-service/app/transfer"
	"github.com/torngkab/grit-api-gateway-service/app/withdraw"
	"github.com/torngkab/grit-api-gateway-service/app/ws"
	"github.com/torngkab/grit-api-gateway-service/config"
	_ "github.com/torngkab/grit-api-gateway-service/docs"

	gatewayMiddleware "github.com/torngkab/grit-api-gateway-service/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewServer(config config.Config) *echo.Echo {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Routes()
	return e
}

func router(e *echo.Echo, config config.Config) {
	e.Validator = &CustomValidator{validator: validator.New()}

	// websocket route
	{
		getBalanceAdapter := getbalance.NewGetBalanceAdapter(config)
		getBalanceService := getbalance.NewGetBalanceService(getBalanceAdapter)
		wsHandler := ws.NewHandler(getBalanceService)
		e.GET("/ws", wsHandler.ServeWebSocket)
	}

	group := e.Group("/api/v1")
	groupAccount := group.Group("/accounts")
	{
		{
			createUserAdapter := createaccount.NewCreateUserAdapter(config)
			createUserService := createaccount.NewCreateUserService(createUserAdapter)
			createUserHandler := createaccount.NewHandler(createUserService)
			groupAccount.POST("", createUserHandler.CreateUser)
		}

		{
			getBalanceAdapter := getbalance.NewGetBalanceAdapter(config)
			getBalanceService := getbalance.NewGetBalanceService(getBalanceAdapter)
			getBalanceHandler := getbalance.NewHandler(getBalanceService)
			groupAccount.GET("/:id/balance", getBalanceHandler.GetBalance, gatewayMiddleware.BasicAuth(config))
		}
	}

	{
		getTransactionsAdapter := gettransactions.NewGetTransactionsAdapter(config)
		getTransactionsService := gettransactions.NewGetTransactionsService(getTransactionsAdapter)
		getTransactionsHandler := gettransactions.NewHandler(getTransactionsService)
		group.GET("/transactions", getTransactionsHandler.GetTransactions, gatewayMiddleware.BasicAuth(config))
	}

	{
		depositAdapter := deposit.NewDepositAdapter(config)
		depositService := deposit.NewDepositService(depositAdapter)
		depositHandler := deposit.NewHandler(depositService)
		group.POST("/deposit", depositHandler.Deposit, gatewayMiddleware.BasicAuth(config))
	}

	{
		withdrawAdapter := withdraw.NewWithdrawAdapter(config)
		withdrawService := withdraw.NewWithdrawService(withdrawAdapter)
		withdrawHandler := withdraw.NewHandler(withdrawService)
		group.POST("/withdraw", withdrawHandler.Withdraw, gatewayMiddleware.BasicAuth(config))
	}

	{
		transferAdapter := transfer.NewTransferAdapter(config)
		transferService := transfer.NewTransferService(transferAdapter)
		transferHandler := transfer.NewHandler(transferService)
		group.POST("/transfer", transferHandler.Transfer, gatewayMiddleware.BasicAuth(config))
	}
}

// @title			Grit Challenge - API Gateway Service
// @version		1.0
// @description	API Gateway Service
// @termsOfService	http://swagger.io/terms/
func main() {
	// get config
	config := config.C("")

	e := NewServer(config)
	router(e, config)

	fmt.Println("Registered Routes:")
	for _, route := range e.Routes() {
		fmt.Printf("Method: %s, Path: %s, Handler: %s\n", route.Method, route.Path, route.Name)
	}

	e.Start(":" + config.Server.Port)
}
