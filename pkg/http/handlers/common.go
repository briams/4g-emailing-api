package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"git.gdpteam.com/4gen/4g-tags-api/pkg/storage"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

// CommonHandler has the param handlers
type CommonHandler struct {
	DB *sql.DB
}

// NewCommonHandler returns a New CommonHandler
func NewCommonHandler(db *sql.DB) *CommonHandler {
	return &CommonHandler{db}
}

// CheckDBHealth godoc
// @Summary returns the time from DB
// @Description returns the time from DB
// @Tags commons
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Success 200 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /health [get]
func (p *CommonHandler) CheckDBHealth(c echo.Context) error {
	mr := utils.ResponseMessage{}

	storageParam := storage.NewMySQLCommon(p.DB)
	now, err := storageParam.CheckDBHealth()
	if err != nil {
		log.Printf("error: Cannot retrieve the now time from DB. Handler common.CheckDBHealth: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"¡Upps! no pudimos crear el registro",
			"para descubrir que sucedio revisa los log del servicio",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	mr.AddMessage(http.StatusOK, "Health Checked, everything is fine.", "")
	mr.Data = struct {
		Now time.Time `json:"now"`
	}{*now}

	return c.JSON(http.StatusOK, mr)
}
