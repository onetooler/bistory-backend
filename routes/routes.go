package routes

import (
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
	setErrorController(e, container)
	setAuthController(e, container)
	setAccountController(e, container)
	setHealthController(e, container)

	setSwagger(container, e)
}

func setErrorController(e *echo.Echo, container container.Container) {
	errorHandler := controller.NewErrorController(container)
	e.HTTPErrorHandler = errorHandler.JSONError
	e.Use(middleware.Recover())
}

func setAuthController(e *echo.Echo, container container.Container) {
	auth := controller.NewAuthController(container)
	e.GET(config.APIAuthLoginStatus, func(c echo.Context) error { return auth.GetLoginStatus(c) })
	e.GET(config.APIAuthLoginAccount, func(c echo.Context) error { return auth.GetLoginAccount(c) })

	if container.GetConfig().Extension.SecurityEnabled {
		e.POST(config.APIAuthLogin, func(c echo.Context) error { return auth.Login(c) })
		e.POST(config.APIAuthLogout, func(c echo.Context) error { return auth.Logout(c) })
		e.POST(config.APIAuthEmailVerificationTokenSend, func(c echo.Context) error { return auth.EmailVerificationTokenSend(c) })
	}
}

func setAccountController(e *echo.Echo, container container.Container) {
	account := controller.NewAccountController(container)
	e.POST(config.APIAccount, func(c echo.Context) error { return account.CreateAccount(c) })
	e.GET(config.APIAccountIdPath, func(c echo.Context) error { return account.GetAccount(c) })
	e.POST(config.APIAccountChangePassword, func(c echo.Context) error { return account.ChangeAccountPassword(c) })
	e.DELETE(config.APIAccountIdPath, func(c echo.Context) error { return account.DeleteAccount(c) })
	e.POST(config.APIAccountFindLoginId, func(c echo.Context) error { return account.FindLoginId(c) })
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
