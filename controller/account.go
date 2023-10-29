package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/service"
)

// AccountController is a controller for managing accounts.
type AccountController interface {
	GetAccount(c echo.Context) error
	CreateAccount(c echo.Context) error
	UpdateAccount(c echo.Context) error
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
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param account_id path int true "Account ID"
// @Success 200 {object} model.Account "Success to fetch data."
// @Failure 400 {string} message "Failed to fetch data."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /accounts/{account_id} [get]
func (controller *accountController) GetAccount(c echo.Context) error {
	account, err := controller.service.FindByLoginId(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, account)
}

// CreateAccount create a new account by http post.
// @Summary Create a new account
// @Description Create a new account
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param data body dto.CreateAccountDto true "a new account data for creating"
// @Success 200 {object} model.Account "Success to create a new account."
// @Failure 400 {string} message "Failed to the registration."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /accounts [post]
func (controller *accountController) CreateAccount(c echo.Context) error {
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

// UpdateAccount update the existing account by http put.
// @Summary Update the existing account
// @Description Update the existing account
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param account_id path int true "Account ID"
// @Param data body dto.UpdatePasswordDto true "the account data for updating"
// @Success 200 {object} model.Account "Success to update the existing account."
// @Failure 400 {string} message "Failed to the update."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /accounts/{account_id} [put]
func (controller *accountController) UpdateAccount(c echo.Context) error {
	dto := dto.NewUpdatePasswordDto()
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, dto)
	}
	account, result := controller.service.UpdateAccountPassword(c.Param("id"), dto)
	if result != nil {
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, account)
}

// DeleteAccount deletes the existing account by http delete.
// @Summary Delete the existing account
// @Description Delete the existing account
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param account_id path int true "Account ID"
// @Success 200 {object} model.Account "Success to delete the existing account."
// @Failure 400 {string} message "Failed to the delete."
// @Failure 401 {boolean} bool "Failed to the authentication. Returns false."
// @Router /accounts/{account_id} [delete]
func (controller *accountController) DeleteAccount(c echo.Context) error {
	account, result := controller.service.DeleteAccount(c.Param("id"))
	if result != nil {
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, account)
}
