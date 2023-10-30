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
	if !controller.container.GetSession().HasAuthorizationTo(accountId, model.AuthorityUser) {
		return c.JSON(http.StatusForbidden, false)
	}

	account, err := controller.service.GetAccount(uint(accountId))
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
	if controller.container.GetSession().GetAccount() != nil {
		return c.String(http.StatusBadRequest, "this account is already logged in")
	}
	dto := dto.NewCreateAccountDto()
	if err := c.Bind(dto); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	account, err := controller.service.CreateAccount(dto)
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
	if !controller.container.GetSession().HasAuthorizationTo(accountId, model.AuthorityUser) {
		return c.JSON(http.StatusForbidden, false)
	}

	dto := dto.NewChangeAccountPasswordDto()
	if err := c.Bind(dto); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	account, err := controller.service.ChangeAccountPassword(accountId, dto)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_ = controller.container.GetSession().Logout()
	return c.JSON(http.StatusOK, account)
}

// DeleteAccount deletes the existing account by http delete.
// @Summary Delete the existing account
// @Description Delete the existing account
// @Tags Account
// @Accept  json
// @Produce  json
// @Param accountId path int true "Account ID"
// @Success 200 {object} model.Account "Success to delete the existing account."
// @Failure 400 {string} message "Failed to the delete."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /account/{accountId} [delete]
func (controller *accountController) DeleteAccount(c echo.Context) error {
	accountId := util.ConvertToUint(c.Param(config.APIAccountIdParam))
	if accountId == 0 {
		return c.String(http.StatusBadRequest, "failed to parse id")
	}
	if !controller.container.GetSession().HasAuthorizationTo(accountId, model.AuthorityUser) {
		return c.JSON(http.StatusForbidden, false)
	}

	account, result := controller.service.DeleteAccount(accountId)
	if result != nil {
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, account)
}
