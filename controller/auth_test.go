package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/onetooler/bistory-backend/util"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginStatus_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.GET(config.APIAuthLoginStatus, func(c echo.Context) error { return auth.GetLoginStatus(c) })

	req := httptest.NewRequest("GET", config.APIAuthLoginStatus, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, "true", rec.Body.String())
}

func TestGetLoginAccount_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.GET(config.APIAuthLoginAccount, func(c echo.Context) error { return auth.GetLoginAccount(c) })

	req := httptest.NewRequest("GET", config.APIAuthLoginAccount, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLogin_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogin, func(c echo.Context) error { return auth.Login(c) })

	param := createLoginSuccessAccount()
	req := testutil.NewJSONRequest("POST", config.APIAuthLogin, param)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, testutil.GetCookie(rec, "GSESSION"))
}

func TestLogin_AuthenticationFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogin, func(c echo.Context) error { return auth.Login(c) })

	param := createLoginFailureAccount()
	req := testutil.NewJSONRequest("POST", config.APIAuthLogin, param)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Empty(t, testutil.GetCookie(rec, "GSESSION"))
}

func TestLogout_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	auth := NewAuthController(container)
	router.POST(config.APIAuthLogout, func(c echo.Context) error { return auth.Logout(c) })

	req := testutil.NewJSONRequest("POST", config.APIAuthLogout, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, testutil.GetCookie(rec, "GSESSION"))
}

func TestEmailVerificationTokenSend_Success(t *testing.T) {
	mailServer := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
		PortNumber:        testutil.TestEmailServerPort,
	})
	err := mailServer.Start()
	assert.Nil(t, err)
	defer util.Check(mailServer.Stop)

	router, container := testutil.PrepareForControllerTest(true)

	auth := NewAuthController(container)
	router.POST(config.APIAuthEmailVerificationTokenSend, func(c echo.Context) error { return auth.EmailVerificationTokenSend(c) })

	dto := dto.EmailVerificationTokenSendDto{
		Email: "newTest@example.com",
	}

	req := testutil.NewJSONRequest("POST", config.APIAuthEmailVerificationTokenSend, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, testutil.GetCookie(rec, "GSESSION"))
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
