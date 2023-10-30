package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockService struct {
	createAccount         func(*dto.CreateAccountDto) (*model.Account, error)
	changeAccountPassword func(uint, *dto.ChangeAccountPasswordDto) (*model.Account, error)
	deleteAccount         func(uint) (bool, error)
	getAccount            func(uint) (*model.Account, error)
}

func (m *mockService) CreateAccount(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
	return m.createAccount(createAccountDto)
}

func (m *mockService) ChangeAccountPassword(id uint, UpdatePasswordDto *dto.ChangeAccountPasswordDto) (*model.Account, error) {
	return m.changeAccountPassword(id, UpdatePasswordDto)
}

func (m *mockService) DeleteAccount(id uint) (bool, error) {
	return m.deleteAccount(id)
}

func (m *mockService) GetAccount(id uint) (*model.Account, error) {
	return m.getAccount(id)
}

func TestCreateAccount_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return &model.Account{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					LoginId:   createAccountDto.LoginId,
					Email:     createAccountDto.Email,
					Password:  "hashed" + createAccountDto.Password,
					Authority: model.AuthorityUser,
				}, nil
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error { return account.CreateAccount(c) })

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTestTest",
	}
	req := test.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body := model.Account{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Nil(t, err)
	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body.ID)
	assert.Equal(t, model.AuthorityUser, body.Authority)
	assert.Equal(t, dto.LoginId, body.LoginId)
	assert.Equal(t, dto.Email, body.Email)
	assert.NotEmpty(t, body.CreatedAt)
	assert.Empty(t, body.Password)
}

func TestCreateAccount_WrongPasswordFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return nil, fmt.Errorf("password must be at least 8 characters")
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error { return account.CreateAccount(c) })

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTest",
	}
	req := test.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAccount_DuplicatedLoginIdFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return nil, fmt.Errorf("duplicated LoginId")
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error { return account.CreateAccount(c) })

	dto := dto.CreateAccountDto{
		LoginId:  "test",
		Email:    "newTest@example.com",
		Password: "newTestTest",
	}
	req := test.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAccount_DuplicatedEmailFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return nil, fmt.Errorf("duplicated Email")
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error { return account.CreateAccount(c) })

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "test@example.com",
		Password: "newTestTest",
	}
	req := test.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAccount_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			getAccount: func(accountId uint) (*model.Account, error) {
				return &testAccount, nil
			},
		},
	}
	router.GET(config.APIAccountIdPath, func(c echo.Context) error {
		login(container, testAccount)
		return account.GetAccount(c)
	})

	req := test.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body := model.Account{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Nil(t, err)
	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body.ID)
	assert.Equal(t, model.AuthorityUser, body.Authority)
	assert.Equal(t, "newTest", body.LoginId)
	assert.Equal(t, "newTest@example.com", body.Email)
	assert.NotEmpty(t, body.CreatedAt)
	assert.Empty(t, body.Password)
}

func TestGetAccount_NoLoginFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			getAccount: func(accountId uint) (*model.Account, error) {
				return &testAccount, nil
			},
		},
	}
	router.GET(config.APIAccountIdPath, func(c echo.Context) error {
		return account.GetAccount(c)
	})

	req := test.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestGetAccount_NoAuthorizationFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			getAccount: func(accountId uint) (*model.Account, error) {
				return &testAccount, nil
			},
		},
	}
	router.GET(config.APIAccountIdPath, func(c echo.Context) error {
		login(container, testAccount)
		return account.GetAccount(c)
	})

	req := test.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID+1), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestChangeAccountPassword_Success(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			changeAccountPassword: func(accountId uint, changeAccountPasswordDto *dto.ChangeAccountPasswordDto) (*model.Account, error) {
				return &testAccount, nil
			},
		},
	}
	router.POST(config.APIAccountChangePassword, func(c echo.Context) error {
		login(container, testAccount)
		return account.ChangeAccountPassword(c)
	})

	req := test.NewJSONRequest(http.MethodPost, strings.Replace(config.APIAccountChangePassword, ":"+config.APIAccountIdParam, strconv.Itoa(int(testAccount.ID)), 1), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body := model.Account{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Nil(t, err)
	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body.ID)
	assert.Equal(t, model.AuthorityUser, body.Authority)
	assert.Equal(t, "newTest", body.LoginId)
	assert.Equal(t, "newTest@example.com", body.Email)
	assert.NotEmpty(t, body.CreatedAt)
	assert.Empty(t, body.Password)
}

func TestChangeAccountPassword_NoAuthorizationFailure(t *testing.T) {
	router, container := test.PrepareForControllerTest(true)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			changeAccountPassword: func(accountId uint, changeAccountPasswordDto *dto.ChangeAccountPasswordDto) (*model.Account, error) {
				return &testAccount, nil
			},
		},
	}
	router.POST(config.APIAccountChangePassword, func(c echo.Context) error {
		return account.ChangeAccountPassword(c)
	})

	req := test.NewJSONRequest(http.MethodPost, strings.Replace(config.APIAccountChangePassword, ":"+config.APIAccountIdParam, strconv.Itoa(int(testAccount.ID)), 1), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func login(container container.Container, account model.Account) {
	_ = container.GetSession().SetAccount(&account)
	_ = container.GetSession().Save()
}

func newTestUserAccount() model.Account {
	return model.Account{
		Model: gorm.Model{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		LoginId:   "newTest",
		Email:     "newTest" + "@example.com",
		Authority: model.AuthorityUser,
	}
}
