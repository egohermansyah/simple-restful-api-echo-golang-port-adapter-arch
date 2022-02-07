package router

import (
	"net/http"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/role"
	"simple-restful-api-echo-golang-port-adapter-archhandlers/user"

	"github.com/labstack/echo/v4"
)

const LOG_IDENTIFIER = "API_ROUTER"

type Handler struct {
	Role role.Handler
	User user.Handler
}

func RegisterPath(
	e *echo.Echo,
	handler Handler,
) {
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	roleV1 := e.Group("api/v1/role")
	roleV1.GET("", handler.Role.List)
	roleV1.GET("/:id", handler.Role.FindById)
	roleV1.PUT("/:id", handler.Role.UpdateById)

	userV1 := e.Group("api/v1/user")
	userV1.GET("", handler.User.List)
	userV1.POST("", handler.User.Create)
	userV1.GET("/:id", handler.User.FindById)
	userV1.PUT("/:id", handler.User.UpdateById)
	userV1.DELETE("/:id", handler.User.DeleteById)
}
