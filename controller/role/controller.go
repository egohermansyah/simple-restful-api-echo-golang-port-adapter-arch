package role

import (
	"golang-vscode-setup/controller/role/defined"
	"golang-vscode-setup/controller/util/queryparams"
	"golang-vscode-setup/controller/util/response"
	"golang-vscode-setup/service/role"
	serviceDefined "golang-vscode-setup/service/role/defined"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service role.IService
}

func NewController(service role.IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller Controller) List(c echo.Context) error {
	queryParams := c.QueryParams()
	cleanQueryParams := queryparams.QueryParamsCleaner(queryParams)
	result, err := controller.service.List(cleanQueryParams.QueryParams, cleanQueryParams.PerPage, cleanQueryParams.Offset)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	var results []defined.DefaultResponse
	for _, data := range result {
		results = append(results, *defined.NewDefaultResponse(data))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], results))
}

func (controller Controller) FindById(c echo.Context) error {
	id := c.Param("id")
	result, err := controller.service.FindById(id)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], defined.NewDefaultResponse(result)))
}

func (controller Controller) UpdateById(c echo.Context) error {
	id := c.Param("id")
	bodyRequest := new(defined.InsertRequest)
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
	data := serviceDefined.Role{Name: bodyRequest.Name, Desc: bodyRequest.Desc}
	result, err := controller.service.UpdateById(id, data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], defined.NewDefaultResponse(result)))
}
