package controller

import (
	"golang-vscode-setup/controller/role"
	"golang-vscode-setup/controller/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

const LOG_IDENTIFIER = "API_ROUTER"

type Controllers struct {
	Role role.Controller
	User user.Controller
}

func RegisterPath(
	e *echo.Echo,
	controllers Controllers,
) {
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	roleV1 := e.Group("api/v1/role")
	roleV1.GET("", controllers.Role.List)
	roleV1.GET("/:id", controllers.Role.FindById)
	roleV1.PUT("/:id", controllers.Role.UpdateById)

	userV1 := e.Group("api/v1/user")
	userV1.GET("", controllers.User.List)
	userV1.POST("", controllers.User.Create)
	userV1.GET("/:id", controllers.User.FindById)
	userV1.PUT("/:id", controllers.User.UpdateById)
	userV1.DELETE("/:id", controllers.User.DeleteById)
}
