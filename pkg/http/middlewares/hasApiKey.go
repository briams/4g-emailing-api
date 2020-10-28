package middlewares

import (
	"log"
	"net/http"
	"os"

	"git.gdpteam.com/4gen/4g-tags-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

// HasAPIKeyHeader validate the API_KEY header
func HasAPIKeyHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		msgResponse := utils.ResponseMessage{}

		APIKEY := os.Getenv("GENERAL_API_KEY")
		apiKeyVal := c.Request().Header.Get("API_KEY")

		if apiKeyVal == "" {
			log.Println("API_KEY: Missing Header")
			msgResponse.AddError(
				http.StatusForbidden,
				"API_KEY header is missing",
				"",
			)
			return c.JSON(http.StatusBadRequest, msgResponse)
		}

		if apiKeyVal != APIKEY {
			log.Println("API_KEY: Header value does not match the API_KEY value")
			msgResponse.AddError(
				http.StatusForbidden,
				"Header value does not match the API_KEY value",
				"",
			)
			return c.JSON(http.StatusBadRequest, msgResponse)
		}

		return next(c)
	}
}
