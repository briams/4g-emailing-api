package routes

import (
	"database/sql"

	"git.gdpteam.com/4gen/4g-tags-api/db/rds"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/http/handlers"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/http/middlewares"
	"github.com/labstack/echo/v4"
)

// ParamRoutes registers param routes
func ParamRoutes(e *echo.Echo, db *sql.DB, rdb *rds.Rds) {
	r := e.Group("/api/v1/models")
	r.Use(middlewares.HasAPIKeyHeader)

	h := handlers.NewTagHandler(db, rdb)
	r.POST("", h.Create)
	r.GET("", h.GetAll)
	r.GET("/:id", h.GetByID)
	// r.GET("/list", h.GetByIDs)
	r.PUT("/:id", h.Update)
}
