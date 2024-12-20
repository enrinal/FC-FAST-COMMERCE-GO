package main

import (
	"database/sql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"order-service/internal/api"
	"order-service/internal/config"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/internal/sharding"
	"time"
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/order-db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	//db, err := connectDB()
	//if err != nil {
	//	panic(err)
	//}

	db1, err := sql.Open("mysql", "root:@tcp(localhost:3306)/order-db-1")
	if err != nil {
		panic(err)
	}

	db2, err := sql.Open("mysql", "root:@tcp(localhost:3306)/order-db-2")
	if err != nil {
		panic(err)
	}

	db3, err := sql.Open("mysql", "root:@tcp(localhost:3306)/order-db-3")
	if err != nil {
		panic(err)
	}

	kafkaWriter := config.NewKafkaWriter("order-topic")

	router := sharding.NewShardRouter(3)

	orderRepo := repository.NewOrderRepository([]*sql.DB{db1, db2, db3}, router)
	orderService := service.NewOrderService(*orderRepo, "http://localhost:8081", "http://localhost:8083", kafkaWriter)
	orderHandler := api.NewOrderHandler(*orderService)

	e := echo.New()

	limiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(1),
				Burst:     3,
				ExpiresIn: 3 * time.Minute,
			}),
		IdentifierExtractor: func(context echo.Context) (string, error) {
			// for local
			return context.Request().RemoteAddr, nil
			// for production
			// return context.Request().Header.Get(echo.HeaderXRealIP), nil
			// return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{"error": "rate limit exceeded"})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(429, map[string]string{"error": "rate limit exceeded"})
		},
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.JWT([]byte("secret")))
	e.Use(middleware.RateLimiterWithConfig(limiterConfig))

	e.POST("/orders", orderHandler.CreateOrder)
	e.PUT("/orders", orderHandler.UpdateOrder)
	e.DELETE("/orders/:id", orderHandler.CancelOrder)

	e.Logger.Fatal(e.Start(":8082"))
}
