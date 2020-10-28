package routes

import (
	"database/sql"

	"github.com/briams/4g-emailing-api/pkg/http/handlers"
	"github.com/labstack/echo/v4"
)

// CommonRoutes registers common routes
func CommonRoutes(e *echo.Echo, db *sql.DB) {
	r := e.Group("/api/v1/health")

	h := handlers.NewCommonHandler(db)
	r.GET("", h.CheckDBHealth)
}
