package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/test"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginStatus_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.GET(config.APIAuthLoginStatus, func(c echo.Context) error { return auth.GetLoginStatus(c) })

	req := httptest.NewRequest("GET", config.APIAuthLoginStatus, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, "true", rec.Body.String())
}

func TestGetLoginAccount_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.GET(config.APIAuthLoginAccount, func(c echo.Context) error { return auth.GetLoginAccount(c) })

	req := httptest.NewRequest("GET", config.APIAuthLoginAccount, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLogin_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogin, func(c echo.Context) error { return auth.Login(c) })

	param := createLoginSuccessAccount()
	req := test.NewJSONRequest("POST", config.APIAuthLogin, param)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, test.GetCookie(rec, "GSESSION"))
}

func TestLogin_AuthenticationFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogin, func(c echo.Context) error { return auth.Login(c) })

	param := createLoginFailureAccount()
	req := test.NewJSONRequest("POST", config.APIAuthLogin, param)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Empty(t, test.GetCookie(rec, "GSESSION"))
}

func TestLogout_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogout, func(c echo.Context) error { return auth.Logout(c) })

	req := test.NewJSONRequest("POST", config.APIAuthLogout, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, test.GetCookie(rec, "GSESSION"))
}

func createLoginSuccessAccount() *dto.LoginDto {
	return &dto.LoginDto{
		LoginId:  "test",
		Email:    "test@example.com",
		Password: "test",
	}
}

func createLoginFailureAccount() *dto.LoginDto {
	return &dto.LoginDto{
		LoginId:  "test",
		Email:    "test@example.com",
		Password: "abcde",
	}
}
