package main

import (
	"log"
	"time"

	_ "github.com/briams/4g-emailing-api/docs"

	"github.com/briams/4g-emailing-api/config"
	"github.com/briams/4g-emailing-api/pkg/http/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func initSetup() *echo.Echo {
	e := echo.New()
	config.EchoServerTimeout(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.EchoAllowOrigins(),
		AllowCredentials: true,
		MaxAge:           int(time.Minute) * 10,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "API_KEY"},
	}))

	db := newConnection()
	// rdb := newRedisClient()
	// refreshData(rdb)

	// App Routes
	routes.ModelRoutes(e, db)
	routes.CommonRoutes(e, db)

	// Swagger Docs
	log.Println("Swagger Docs are in /swagger/index.html")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
