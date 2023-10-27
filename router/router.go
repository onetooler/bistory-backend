package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/controller"

	_ "github.com/onetooler/bistory-backend/docs" // for using echo-swagger
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Init initialize the routing of this application.
func Init(e *echo.Echo, container container.Container) {
	setCORSConfig(e, container)

	setErrorController(e, container)
	setAccountController(e, container)
	setHealthController(e, container)

	setSwagger(container, e)
}

func setCORSConfig(e *echo.Echo, container container.Container) {
	if container.GetConfig().Extension.CorsEnabled {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials:                         true,
			UnsafeWildcardOriginWithAllowCredentials: true,
			AllowOrigins:                             []string{"*"},
			AllowHeaders: []string{
				echo.HeaderAccessControlAllowHeaders,
				echo.HeaderContentType,
				echo.HeaderContentLength,
				echo.HeaderAcceptEncoding,
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			MaxAge: 86400,
		}))
	}
}

func setErrorController(e *echo.Echo, container container.Container) {
	errorHandler := controller.NewErrorController(container)
	e.HTTPErrorHandler = errorHandler.JSONError
	e.Use(middleware.Recover())
}

func setAccountController(e *echo.Echo, container container.Container) {
	account := controller.NewAccountController(container)
	e.GET(config.APIAccountLoginStatus, func(c echo.Context) error { return account.GetLoginStatus(c) })
	e.GET(config.APIAccountLoginAccount, func(c echo.Context) error { return account.GetLoginAccount(c) })

	if container.GetConfig().Extension.SecurityEnabled {
		e.POST(config.APIAccountLogin, func(c echo.Context) error { return account.Login(c) })
		e.POST(config.APIAccountLogout, func(c echo.Context) error { return account.Logout(c) })
	}
}

func setHealthController(e *echo.Echo, container container.Container) {
	health := controller.NewHealthController(container)
	e.GET(config.APIHealth, func(c echo.Context) error { return health.GetHealthCheck(c) })
}

func setSwagger(container container.Container, e *echo.Echo) {
	if container.GetConfig().Swagger.Enabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
