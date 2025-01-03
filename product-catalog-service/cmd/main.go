package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"product-catalog-service/internal/api"
	consumer2 "product-catalog-service/internal/consumer"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"
	"time"
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/product-db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize product service
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(*productRepo, rdb)
	productHandler := api.NewProductHandler(*productService)

	// consumer
	consumer := consumer2.NewConsumer(productService)
	go consumer.StartKafkaConsumer()

	// Initialize echo
	e := echo.New()

	config := middleware.RateLimiterConfig{
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
	e.Use(middleware.RateLimiterWithConfig(config))

	// Routes
	e.GET("/products/:id/stock", productHandler.GetProductStock)
	e.POST("/products/reserve", productHandler.ReserveProductStock)
	e.POST("/products/release", productHandler.ReleaseProductStock)
	e.GET("/products/warmup-cache", productHandler.PreWarmupCache)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
