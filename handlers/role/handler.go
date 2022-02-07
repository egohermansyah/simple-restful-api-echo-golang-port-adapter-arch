package role

import (
	"net/http"
	domain "simple-restful-api-echo-golang-port-adapter-archcore/domains/role"
	port "simple-restful-api-echo-golang-port-adapter-archcore/ports/role"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/util/queryparams"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/util/response"
	"time"

	"github.com/labstack/echo/v4"
)

type InsertRequest struct {
	Name string `json:"name" validate:"required,min=1"`
	Desc string `json:"desc"`
}

type DefaultResponse struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"updated"`
}

func NewDefaultResponse(data *domain.Role) *DefaultResponse {
	return &DefaultResponse{
		Id:       data.Id,
		Name:     data.Name,
		Desc:     data.Desc,
		Created:  data.Created,
		Modified: data.Modified,
	}
}

type Handler struct {
	service port.Service
}

func New(service port.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler Handler) List(c echo.Context) error {
	queryParams := c.QueryParams()
	cleanQueryParams := queryparams.QueryParamsCleaner(queryParams)
	result, err := handler.service.List(cleanQueryParams.QueryParams, cleanQueryParams.PerPage, cleanQueryParams.Offset)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	var results []DefaultResponse
	for _, data := range result {
		results = append(results, *NewDefaultResponse(data))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], results))
}

func (handler Handler) FindById(c echo.Context) error {
	id := c.Param("id")
	result, err := handler.service.FindById(id)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], NewDefaultResponse(result)))
}

func (handler Handler) UpdateById(c echo.Context) error {
	id := c.Param("id")
	bodyRequest := new(InsertRequest)
	if err := c.Bind(bodyRequest); err != nil {
		result := response.Map["badRequest"]
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	if err := c.Validate(bodyRequest); err != nil {
		errors := response.BuildErrorBodyRequestValidator(err)
		result := response.Map["badRequest"]
		result.Errors = errors
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	data := domain.Role{Name: bodyRequest.Name, Desc: bodyRequest.Desc}
	result, err := handler.service.UpdateById(id, data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], NewDefaultResponse(result)))
}
