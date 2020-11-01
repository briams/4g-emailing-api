package routes

import (
	"database/sql"

	"github.com/briams/4g-emailing-api/pkg/http/handlers"
	"github.com/briams/4g-emailing-api/pkg/http/middlewares"
	"github.com/labstack/echo/v4"
)

// ParamRoutes registers param routes
func ParamRoutes(e *echo.Echo, db *sql.DB) {
	r := e.Group("/api/v1/models")
	r.Use(middlewares.HasAPIKeyHeader)

	h := handlers.NewTagHandler(db)
	r.POST("", h.Create)
	r.GET("", h.GetAll)
	r.GET("/:id", h.GetByID)
	r.PUT("/:id", h.Update)
	r.PATCH("/:id/activate", h.Activate)
	r.PATCH("/:id/deactivate", h.Deactivate)

}
