package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/service"
	"github.com/onetooler/bistory-backend/session"
)

// AuthController is a controller for managing user account.
type AuthController interface {
	GetLoginStatus(c echo.Context) error
	GetLoginAccount(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type authController struct {
	container container.Container
	service   service.AuthService
}

// NewAuthController is constructor.
func NewAuthController(container container.Container) AuthController {
	return &authController{
		container: container,
		service:   service.NewAuthService(container),
	}
}

// GetLoginStatus returns the status of login.
// @Summary Get the login status.
// @Description Get the login status of current logged-in user.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {boolean} bool "The current user have already logged-in. Returns true."
// @Failure 401 {boolean} bool "The current user haven't logged-in yet. Returns false."
// @Router /auth/loginStatus [get]
func (controller *authController) GetLoginStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}

// GetLoginAccount returns the account data of logged in user.
// @Summary Get the account data of logged-in user.
// @Description Get the account data of logged-in user.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Account "Success to fetch the account data. If the security function is disable, it returns disabled message"
// @Failure 401 {boolean} bool "The current user haven't logged-in yet. Returns false."
// @Router /auth/loginAccount [get]
func (controller *authController) GetLoginAccount(c echo.Context) error {
	if !controller.container.GetConfig().Extension.SecurityEnabled {
		return c.JSON(http.StatusOK, "Security is disabled")
	}
	return c.JSON(http.StatusOK, controller.container.GetSession().GetAccount())
}

// Login is the method to login using loginId and password by http post.
// @Summary Login using loginId and password.
// @Description Login using loginId and password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param data body dto.LoginDto true "User name and Password for logged-in."
// @Success 200 {object} model.Account "Success to the authentication."
// @Failure 401 {boolean} bool "Failed to the authentication."
// @Router /auth/login [post]
func (controller *authController) Login(c echo.Context) error {
	dto := dto.NewLoginDto()
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, dto)
	}

	sess := controller.container.GetSession()
	if account := sess.GetAccount(); account != nil {
		return c.JSON(http.StatusOK, account)
	}

	authorized, account := controller.service.AuthenticateByLoginIdAndPassword(dto.LoginId, dto.Password)
	if !authorized {
		return c.JSON(http.StatusForbidden, false)
	}
	_ = sess.Login(
		&session.Account{
			Id:        account.ID,
			LoginId:   account.LoginId,
			Authority: uint(account.Authority),
		},
	)
	return c.JSON(http.StatusOK, account)
}

// Logout is the method to logout by http post.
// @Summary Logout.
// @Description Logout.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200
// @Router /auth/logout [post]
func (controller *authController) Logout(c echo.Context) error {
	_ = controller.container.GetSession().Logout()
	return c.NoContent(http.StatusOK)
}
