package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/briams/4g-emailing-api/pkg/http/validators"
	"github.com/briams/4g-emailing-api/pkg/models/model"
	"github.com/briams/4g-emailing-api/pkg/storage"
	"github.com/briams/4g-emailing-api/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// ModelHandler has the model handlers
type ModelHandler struct {
	DB *sql.DB
}

// NewModelHandler returns a New ModelHandler
func NewModelHandler(db *sql.DB) *ModelHandler {
	return &ModelHandler{db}
}

// Create godoc
// @Summary Create a model
// @Description Create a new model item
// @Models models
// @Accept json
// @Produce json
// @Param model body validators.CreateBody true "New model"
// @Param API_KEY header string required "API_KEY Header"
// @Success 201 {object} utils.ResponseMessage
// @Failure 422 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models [post]
func (p *ModelHandler) Create(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.CreateBody{}
	m := &model.Model{}

	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: model structure is not correct Handler model.Create: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"Must be a valid structure",
			"check the logs for more details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: invalid data. Handler model.Create: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"Invalid Data",
			"check the documentation",
		)
		mr.Data = errors
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	m.ModelID = b.ModelID
	m.Mjml = b.Mjml
	m.Variables = b.Variables
	m.InsUserID = b.InsUserID

	storageModel := storage.NewMySQLModel(p.DB, nil, nil)
	serviceModel := model.NewService(storageModel)
	err = serviceModel.Create(m)
	if errors.Is(err, model.ErrModelAlreadyExist) {
		log.Printf("error: cannot create a Model. Handler model.Create: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			model.ErrModelAlreadyExist.Error(),
			"Check the service logs for details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}
	if err != nil {
		log.Printf("error: cannot create a Model. Handler model.Create: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot create a Model",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ = serviceModel.GetByID(m.ModelID)

	mr.AddMessage(http.StatusCreated, "Model Created Successfully", "")
	mr.Data = m

	return c.JSON(http.StatusCreated, mr)
}

// GetAll godoc
// @Summary Get all models
// @Description Get all the models by defining the fields
// @Models models
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param fields query []string true "Event fields"
// @Success 200 {object} utils.ResponseMessage
// @Failure 400 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models [get]
func (p *ModelHandler) GetAll(c echo.Context) error {
	mr := utils.ResponseMessage{}

	storageModel := storage.NewMySQLModel(p.DB, nil, nil)
	serviceModel := model.NewService(storageModel)
	res, err := serviceModel.GetAll()
	if errors.Is(err, model.ErrFieldsDoesNotExist) {
		log.Printf("error: no se pudo obtener la información. Handler event.GetAll: %v", err)
		mr.AddError(
			http.StatusBadRequest,
			"Cannot retrieve the information. "+err.Error(),
			"Check the service logs for details",
		)
		return c.JSON(http.StatusBadRequest, mr)
	}
	if err != nil {
		log.Printf("error: no se pudo obtener la información. Handler user.GetAll: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot retrieve the information.",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	mr.AddMessage(http.StatusOK, "Models fetched successfully.", "")
	mr.Data = res

	return c.JSON(http.StatusOK, mr)
}

// GetByID godoc
// @Summary Get a model
// @Description Get a model by id
// @Models models
// @Accept json
// @Produce json
// @Param API_KEY header string true "API_KEY Header"
// @Param id path string true "Model ID"
// @Success 200 {object} utils.ResponseMessage
// @Success 422 {object} utils.ResponseMessage
// @Success 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models/{id} [get]
func (p *ModelHandler) GetByID(c echo.Context) error {
	mr := utils.ResponseMessage{}

	id := c.Param("id")

	storageModel := storage.NewMySQLModel(p.DB, nil, nil)
	serviceModel := model.NewService(storageModel)
	res, err := serviceModel.GetByID(id)

	if errors.Is(err, redis.Nil) {
		mr.AddMessage(
			http.StatusNotFound,
			"model not found with id "+id,
			"",
		)
		return c.JSON(http.StatusNotFound, mr)
	}
	if err != nil {
		log.Printf("error: no se pudo obtener los datos solicitados del id: %d. Handler model.GetByID: %v", id, err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot retrieve the model",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	mr.AddMessage(http.StatusOK, "¡listo!", "")
	mr.Data = res

	return c.JSON(http.StatusOK, mr)
}

// Update godoc
// @Summary Update a model
// @Description Update a model item
// @Models models
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param id path string true "model ID"
// @Param model body validators.UpdateBody true "model Updated"
// @Success 200 {object} utils.ResponseMessage
// @Failure 422 {object} utils.ResponseMessage
// @Failure 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models/{id} [put]
func (p *ModelHandler) Update(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.UpdateBody{}
	m := &model.Model{}

	id := c.Param("id")

	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: incorrect structure. Handler model.Update: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"A correct structure must be sended",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: invalid data. Handler model.Update: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"Invalid Data",
			"check the documentation",
		)
		mr.Data = errors
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	m.ModelID = id
	m.Mjml = b.Mjml
	m.Variables = b.Variables
	m.InsUserID = b.SetUserID

	storageAudit := storage.NewMySQLModelAudit(p.DB)

	storageModel := storage.NewMySQLModel(p.DB, storageAudit, nil)
	serviceModel := model.NewService(storageModel)
	err = serviceModel.Update(m)

	if errors.Is(err, sql.ErrNoRows) {
		mr.AddMessage(
			http.StatusNotFound,
			"model not found with id "+id,
			"",
		)
		return c.JSON(http.StatusNotFound, mr)
	}
	if err != nil {
		log.Printf("error: Cannot update the model Handler model.Update: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot update the model",
			"para descubrir que sucedio revisa los log del servicio",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ = serviceModel.GetByID(id)

	mr.AddMessage(http.StatusOK, "Model Updated Successfully.", "")
	mr.Data = m

	return c.JSON(http.StatusOK, mr)
}

// Activate godoc
// @Summary Activate a model
// @Description Activate a model item
// @Models models
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param id path string true "Model ID"
// @Param model body models.ActiveLogBody true "Model Activated"
// @Success 200 {object} utils.ResponseMessage
// @Failure 400 {object} utils.ResponseMessage
// @Failure 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models/{id}/activate [patch]
func (p *ModelHandler) Activate(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.ActiveLogBody{}

	id := c.Param("id")
	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: Incorrect structure. Handler model.Activate: %v", err)
		mr.AddError(
			http.StatusBadRequest,
			"A correct structure must be sended",
			"Check logs on the server for more details",
		)
		return c.JSON(http.StatusBadRequest, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: Invalid data. Handler model.Activate: %v", errors)
		mr.AddError(
			http.StatusBadRequest,
			"Invalid Data",
			"Check the documentation for more details",
		)
		mr.Data = errors
		return c.JSON(http.StatusBadRequest, mr)
	}

	// Enviando datos al Storage(Repository) storage/model
	storageActiveLog := storage.NewMySQLModelActiveLog(p.DB)

	storageModel := storage.NewMySQLModel(p.DB, nil, storageActiveLog)
	serviceModel := model.NewService(storageModel)
	err = serviceModel.Activate(id, b.Reason, b.SetUserID)

	if err == sql.ErrNoRows {
		mr.AddMessage(
			http.StatusNoContent,
			"model not found with id "+id,
			"",
		)
		return c.JSON(http.StatusOK, mr)
	}
	if err != nil {
		log.Printf("error: model not found with id: %s. Handler model.Activate: %v", id, err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot fetch the information",
			"Check logs on the server for more details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ := storageModel.GetByID(id)

	mr.AddMessage(http.StatusOK, "model Activated successfully.", "")
	mr.Data = m

	return c.JSON(http.StatusOK, mr)
}

// Deactivate godoc
// @Summary Deactivate a model
// @Description Deactivate a model item
// @Models models
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param id path string true "Model ID"
// @Param model body models.ActiveLogBody true "Model Deactivate"
// @Success 200 {object} utils.ResponseMessage
// @Failure 400 {object} utils.ResponseMessage
// @Failure 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /models/{id}/deactivate [patch]
func (p *ModelHandler) Deactivate(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.ActiveLogBody{}

	id := c.Param("id")
	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: Incorrect structure. Handler model.Deactivate: %v", err)
		mr.AddError(
			http.StatusBadRequest,
			"A correct structure must be sended",
			"Check logs on the server for more details",
		)
		return c.JSON(http.StatusBadRequest, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: invalid data. Handler model.Deactivate: %v", errors)
		mr.AddError(
			http.StatusBadRequest,
			"Invalid Data",
			"Check the documentation for more details",
		)
		mr.Data = errors
		return c.JSON(http.StatusBadRequest, mr)
	}

	// Enviando datos al Storage(Repository) storage/service
	storageActiveLog := storage.NewMySQLModelActiveLog(p.DB)

	storageModel := storage.NewMySQLModel(p.DB, nil, storageActiveLog)
	serviceModel := model.NewService(storageModel)
	err = serviceModel.Deactivate(id, b.Reason, b.SetUserID)
	if err == sql.ErrNoRows {
		mr.AddMessage(
			http.StatusNoContent,
			"model not found with id "+id,
			"",
		)
		return c.JSON(http.StatusOK, mr)
	}
	if err != nil {
		log.Printf("error: model not found with id: %s. Handler model.Deactivate: %v", id, err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot fetch the information",
			"Check logs on the server for more details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ := storageModel.GetByID(id)

	mr.AddMessage(http.StatusOK, "model Deactivated successfully.", "")
	mr.Data = m

	return c.JSON(http.StatusOK, mr)
}
