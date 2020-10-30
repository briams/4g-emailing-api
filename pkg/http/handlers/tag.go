package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/briams/4g-emailing-api/db/rds"
	"github.com/briams/4g-emailing-api/pkg/http/validators"
	"github.com/briams/4g-emailing-api/pkg/models/tag"
	"github.com/briams/4g-emailing-api/pkg/storage"
	"github.com/briams/4g-emailing-api/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// TagHandler has the tag handlers
type TagHandler struct {
	DB        *sql.DB
	RdsClient *rds.Rds
}

// NewTagHandler returns a New TagHandler
func NewTagHandler(db *sql.DB, rdb *rds.Rds) *TagHandler {
	return &TagHandler{db, rdb}
}

// Create godoc
// @Summary Create a tag
// @Description Create a new tag item
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body validators.CreateBody true "New tag"
// @Param API_KEY header string required "API_KEY Header"
// @Success 201 {object} utils.ResponseMessage
// @Failure 422 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /tags [post]
func (p *TagHandler) Create(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.CreateBody{}
	m := &tag.Model{}

	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: tag structure is not correct Handler tag.Create: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"Must be a valid structure",
			"check the logs for more details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: invalid data. Handler tag.Create: %v", err)
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

	storageTag := storage.NewRedisTag(p.DB, p.RdsClient)
	serviceTag := tag.NewService(storageTag)
	err = serviceTag.Create(m)
	if errors.Is(err, tag.ErrModelAlreadyExist) {
		log.Printf("error: cannot create a Tag. Handler tag.Create: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			tag.ErrModelAlreadyExist.Error(),
			"Check the service logs for details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}
	if err != nil {
		log.Printf("error: cannot create a Tag. Handler tag.Create: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot create a Tag",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ = serviceTag.GetByID(m.ModelID)

	mr.AddMessage(http.StatusCreated, "Tag Created Successfully", "")
	mr.Data = m

	return c.JSON(http.StatusCreated, mr)
}

// GetAll godoc
// @Summary Get all tags
// @Description Get all the tags by defining the fields
// @Tags tags
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param fields query []string true "Event fields"
// @Success 200 {object} utils.ResponseMessage
// @Failure 400 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /tags [get]
func (p *TagHandler) GetAll(c echo.Context) error {
	mr := utils.ResponseMessage{}

	// tagFields := c.QueryParam("fields")
	// tagFieldsSlice := strings.Split(tagFields, ",")
	// if len(tagFields) == 0 {
	// 	log.Println("warning: There is not any event fields")
	// 	mr.AddError(
	// 		http.StatusBadRequest,
	// 		"One event field must be sended at least in the query event: fields",
	// 		"Check the service logs for details",
	// 	)
	// 	return c.JSON(http.StatusBadRequest, mr)
	// }

	storageTag := storage.NewRedisTag(p.DB, p.RdsClient)
	serviceTag := tag.NewService(storageTag)
	// res, err := serviceTag.GetAllWithFields(tagFieldsSlice...)
	res, err := serviceTag.GetAll()
	if errors.Is(err, tag.ErrFieldsDoesNotExist) {
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

	mr.AddMessage(http.StatusOK, "Params fetched successfully.", "")
	mr.Data = res

	return c.JSON(http.StatusOK, mr)
}

// GetByID godoc
// @Summary Get a tag
// @Description Get a tag by id
// @Tags tags
// @Accept json
// @Produce json
// @Param API_KEY header string true "API_KEY Header"
// @Param id path string true "Tag ID"
// @Success 200 {object} utils.ResponseMessage
// @Success 422 {object} utils.ResponseMessage
// @Success 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /tags/{id} [get]
func (p *TagHandler) GetByID(c echo.Context) error {
	mr := utils.ResponseMessage{}

	id := c.Param("id")
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	log.Printf("the ID must be numeric. Handler tag.GetByID: %v", err)
	// 	mr.AddError(
	// 		http.StatusUnprocessableEntity,
	// 		"¡Upps! el id que nos enviaste no es un número entero",
	// 		"",
	// 	)
	// 	return c.JSON(http.StatusUnprocessableEntity, mr)
	// }

	storageTag := storage.NewRedisTag(p.DB, p.RdsClient)
	serviceTag := tag.NewService(storageTag)
	res, err := serviceTag.GetByID(id)

	if errors.Is(err, redis.Nil) {
		mr.AddMessage(
			http.StatusNotFound,
			"tag not found with id "+id,
			"",
		)
		return c.JSON(http.StatusNotFound, mr)
	}
	if err != nil {
		log.Printf("error: no se pudo obtener los datos solicitados del id: %d. Handler tag.GetByID: %v", id, err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot retrieve the tag",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	mr.AddMessage(http.StatusOK, "¡listo!", "")
	mr.Data = res

	return c.JSON(http.StatusOK, mr)
}

// Update godoc
// @Summary Update a tag
// @Description Update a tag item
// @Tags tags
// @Accept json
// @Produce json
// @Param API_KEY header string required "API_KEY Header"
// @Param id path string true "tag ID"
// @Param tag body validators.UpdateBody true "tag Updated"
// @Success 200 {object} utils.ResponseMessage
// @Failure 422 {object} utils.ResponseMessage
// @Failure 404 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /tags/{id} [put]
func (p *TagHandler) Update(c echo.Context) error {
	mr := utils.ResponseMessage{}
	b := &validators.UpdateBody{}
	m := &tag.Model{}

	id := c.Param("id")
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	log.Printf("the ID must be numeric. Handler tag.GetByID: %v", err)
	// 	mr.AddError(
	// 		http.StatusUnprocessableEntity,
	// 		"¡Upps! el id que nos enviaste no es un número entero",
	// 		"",
	// 	)
	// 	return c.JSON(http.StatusUnprocessableEntity, mr)
	// }

	err := c.Bind(b)
	if err != nil {
		log.Printf("warning: incorrect structure. Handler tag.Update: %v", err)
		mr.AddError(
			http.StatusUnprocessableEntity,
			"A correct structure must be sended",
			"Check the service logs for details",
		)
		return c.JSON(http.StatusUnprocessableEntity, mr)
	}

	if errors := b.Validate(); errors != nil {
		log.Printf("warning: invalid data. Handler tag.Update: %v", err)
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

	storageTag := storage.NewRedisTag(p.DB, p.RdsClient)
	serviceTag := tag.NewService(storageTag)
	err = serviceTag.Update(m)
	if errors.Is(err, sql.ErrNoRows) {
		mr.AddMessage(
			http.StatusNotFound,
			"tag not found with id "+id,
			"",
		)
		return c.JSON(http.StatusNotFound, mr)
	}
	if err != nil {
		log.Printf("error: Cannot update the tag Handler tag.Update: %v", err)
		mr.AddError(
			http.StatusInternalServerError,
			"Cannot update the tag",
			"para descubrir que sucedio revisa los log del servicio",
		)
		return c.JSON(http.StatusInternalServerError, mr)
	}

	m, _ = serviceTag.GetByID(id)

	mr.AddMessage(http.StatusOK, "Param Updated Successfully.", "")
	mr.Data = m

	return c.JSON(http.StatusOK, mr)
}
