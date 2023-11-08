package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/service"
	"github.com/onetooler/bistory-backend/util"
)

// AccountController is a controller for managing accounts.
type AccountController interface {
	GetAccount(c echo.Context) error
	CreateAccount(c echo.Context) error
	ChangeAccountPassword(c echo.Context) error
	DeleteAccount(c echo.Context) error
	FindLoginId(c echo.Context) error
}

type accountController struct {
	container container.Container
	service   service.AccountService
}

// NewAccountController is constructor.
func NewAccountController(container container.Container) AccountController {
	return &accountController{container: container, service: service.NewAccountService(container)}
}

// GetAccount returns one record matched account's id.
// @Summary Get a account
// @Description Get a account
// @Tags Account
// @Accept  json
// @Produce  json
// @Param accountId path int true "Account ID"
// @Success 200 {object} model.Account "Success to fetch data."
// @Failure 400 {string} message "Failed to fetch data."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account/{accountId} [get]
func (controller *accountController) GetAccount(c echo.Context) error {
	accountId := util.ConvertToUint(c.Param(config.APIAccountIdParam))
	if accountId == 0 {
		return c.String(http.StatusBadRequest, "failed to parse id")
	}
	if !controller.container.GetSession().HasAuthorizationTo(c, accountId, uint(model.AuthorityUser)) {
		return c.JSON(http.StatusForbidden, false)
	}

	account, err := controller.service.GetAccount(accountId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, account)
}

// CreateAccount create a new account by http post.
// @Summary Create a new account
// @Description Create a new account
// @Tags Account
// @Accept  json
// @Produce  json
// @Param data body dto.CreateAccountDto true "a new account data for creating"
// @Success 200 {object} model.Account "Success to create a new account."
// @Failure 400 {string} message "Failed to the registration."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account [post]
func (controller *accountController) CreateAccount(c echo.Context) error {
	if controller.container.GetSession().GetAccount(c) != nil {
		return c.String(http.StatusBadRequest, "this account is already logged in")
	}
	data := dto.NewCreateAccountDto()
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	account, err := controller.service.CreateAccount(data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, account)
}

// ChangeAccountPassword change account password by http post.
// @Summary Change account password
// @Description Change account password
// @Tags Account
// @Accept  json
// @Produce  json
// @Param accountId path int true "Account ID"
// @Param data body dto.ChangeAccountPasswordDto true "the account password data for updating"
// @Success 200 {object} model.Account "Success to change the account password."
// @Failure 400 {string} message "Failed to the update."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account/{accountId}/ [post]
func (controller *accountController) ChangeAccountPassword(c echo.Context) error {
	accountId := util.ConvertToUint(c.Param(config.APIAccountIdParam))
	if accountId == 0 {
		return c.String(http.StatusBadRequest, "failed to parse id")
	}
	if !controller.container.GetSession().HasAuthorizationTo(c, accountId, uint(model.AuthorityUser)) {
		return c.JSON(http.StatusForbidden, false)
	}

	data := dto.NewChangeAccountPasswordDto()
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	account, err := controller.service.ChangeAccountPassword(accountId, data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = controller.container.GetSession().Logout(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

// DeleteAccount deletes the existing account by http delete.
// @Summary Delete the existing account
// @Description Delete the existing account
// @Tags Account
// @Accept  json
// @Produce  json
// @Param accountId path int true "Account ID"
// @Param data body dto.DeleteAccountDto true "the account password data for updating"
// @Success 200 {boolean} bool "Success to delete the existing account."
// @Failure 400 {string} message "Failed to the delete."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account/{accountId} [delete]
func (controller *accountController) DeleteAccount(c echo.Context) error {
	accountId := util.ConvertToUint(c.Param(config.APIAccountIdParam))
	if accountId == 0 {
		return c.String(http.StatusBadRequest, "failed to parse id")
	}
	if !controller.container.GetSession().HasAuthorizationTo(c, accountId, uint(model.AuthorityUser)) {
		return c.JSON(http.StatusForbidden, false)
	}

	data := dto.NewDeleteAccountDto()
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := controller.service.DeleteAccount(accountId, data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = controller.container.GetSession().Logout(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// FindLoginId send email that contains account's login id to account's email address.
// @Summary Find LoginId By Email
// @Description Find LoginId By Email
// @Tags Account
// @Accept  json
// @Produce  json
// @Param email body dto.FindLoginIdDto true "Account Email"
// @Success 200 {boolean} bool "Success to send email."
// @Failure 400 {string} message "Failed to send email."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account/find-login-id [post]
func (controller *accountController) FindLoginId(c echo.Context) error {
	if controller.container.GetSession().GetAccount(c) != nil {
		return c.String(http.StatusBadRequest, "this account is already logged in")
	}

	data := dto.NewFindLoginIdDto()
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err := controller.service.FindAccountByEmail(data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, true)
}
