package user

import (
	"net/http"
	domain "simple-restful-api-echo-golang-port-adapter-archcore/domains/user"
	port "simple-restful-api-echo-golang-port-adapter-archcore/ports/user"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/util/queryparams"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/util/response"
	"time"

	"github.com/labstack/echo/v4"
)

type InsertRequest struct {
	RoleId      uint   `json:"role_id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateByIdRequest struct {
	RoleId      uint   `json:"role_id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	IsLogin     bool   `json:"is_login"`
	IsActive    bool   `json:"is_active"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	RoleId       uint      `json:"role_id"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	LoginAttempt uint8     `json:"login_attempt"`
	IsLogin      bool      `json:"is_login"`
	IsActive     bool      `json:"is_active"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}

func NewDefaultResponse(data *domain.User) *DefaultResponse {
	return &DefaultResponse{
		Id:           data.Id,
		RoleId:       data.RoleId,
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		LoginAttempt: data.LoginAttempt,
		IsLogin:      data.IsLogin,
		IsActive:     data.IsActive,
		Created:      data.Created,
		Modified:     data.Modified,
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

func (handler Handler) Create(c echo.Context) error {
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
	data := domain.User{
		RoleId:      bodyRequest.RoleId,
		Email:       bodyRequest.Email,
		Password:    bodyRequest.Password,
		PhoneNumber: bodyRequest.PhoneNumber,
	}
	result, err := handler.service.Create(data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusCreated, response.NewResponse("", response.Map["created"], NewDefaultResponse(result)))
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
	bodyRequest := new(UpdateByIdRequest)
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
	data := domain.User{
		RoleId:      bodyRequest.RoleId,
		Email:       bodyRequest.Email,
		Password:    bodyRequest.Password,
		PhoneNumber: bodyRequest.PhoneNumber,
		IsLogin:     bodyRequest.IsLogin,
		IsActive:    bodyRequest.IsActive,
	}
	result, err := handler.service.UpdateById(id, data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], NewDefaultResponse(result)))
}

func (handler Handler) DeleteById(c echo.Context) error {
	id := c.Param("id")
	err := handler.service.DeleteById(id)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["deleted"], nil))
}
