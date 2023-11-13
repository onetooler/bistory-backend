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
	"github.com/onetooler/bistory-backend/infrastructure"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockService struct {
	createAccount         func(*dto.CreateAccountDto) (*model.Account, error)
	changeAccountPassword func(uint, *dto.ChangeAccountPasswordDto) (*model.Account, error)
	deleteAccount         func(uint, *dto.DeleteAccountDto) error
	getAccount            func(uint) (*model.Account, error)
	findAccountByEmail    func(*dto.FindLoginIdDto) error
}

func (m *mockService) CreateAccount(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
	return m.createAccount(createAccountDto)
}

func (m *mockService) ChangeAccountPassword(id uint, UpdatePasswordDto *dto.ChangeAccountPasswordDto) (*model.Account, error) {
	return m.changeAccountPassword(id, UpdatePasswordDto)
}

func (m *mockService) DeleteAccount(id uint, dto *dto.DeleteAccountDto) error {
	return m.deleteAccount(id, dto)
}

func (m *mockService) GetAccount(id uint) (*model.Account, error) {
	return m.getAccount(id)
}

func (m *mockService) FindAccountByEmail(dto *dto.FindLoginIdDto) error {
	return m.findAccountByEmail(dto)
}

func TestCreateAccount_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

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
	router.POST(config.APIAccount, func(c echo.Context) error {
		_ = container.GetSession().SetEmailVerification(c,
			&infrastructure.EmailVerification{Email: "newTest@example.com", VerifiedAt: time.Now()},
		)
		return account.CreateAccount(c)
	})

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTestTest",
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
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
	router, container := testutil.PrepareForControllerTest(false)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return nil, fmt.Errorf("password must be at least 8 characters")
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error {
		_ = container.GetSession().SetEmailVerification(c,
			&infrastructure.EmailVerification{Email: "newTest@example.com", VerifiedAt: time.Now()},
		)
		return account.CreateAccount(c)
	})

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTest",
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAccount_DuplicatedUniqueValueFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	account := accountController{
		container,
		&mockService{
			createAccount: func(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
				return nil, fmt.Errorf("duplicated Email or LoginId")
			},
		},
	}
	router.POST(config.APIAccount, func(c echo.Context) error {
		_ = container.GetSession().SetEmailVerification(c,
			&infrastructure.EmailVerification{Email: "newTest@example.com", VerifiedAt: time.Now()},
		)
		return account.CreateAccount(c)
	})

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "test@example.com",
		Password: "newTestTest",
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAccount_NoEmailVerificationFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

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
	router.POST(config.APIAccount, func(c echo.Context) error {
		_ = container.GetSession().SetEmailVerification(c,
			&infrastructure.EmailVerification{Email: "newTest@example.com", VerifiedAt: time.Time{}},
		)
		return account.CreateAccount(c)
	})

	dto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "test@example.com",
		Password: "newTestTest",
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccount, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAccount_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

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
		login(container, c, testAccount)
		return account.GetAccount(c)
	})

	req := testutil.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID), nil)
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
	router, container := testutil.PrepareForControllerTest(false)

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

	req := testutil.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestGetAccount_NoAuthorizationFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

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
		login(container, c, testAccount)
		return account.GetAccount(c)
	})

	req := testutil.NewJSONRequest(http.MethodGet, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID+1), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestChangeAccountPassword_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

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
		login(container, c, testAccount)
		return account.ChangeAccountPassword(c)
	})

	req := testutil.NewJSONRequest(http.MethodPost, strings.Replace(config.APIAccountChangePassword, ":"+config.APIAccountIdParam, strconv.Itoa(int(testAccount.ID)), 1), nil)
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
	router, container := testutil.PrepareForControllerTest(false)

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

	req := testutil.NewJSONRequest(http.MethodPost, strings.Replace(config.APIAccountChangePassword, ":"+config.APIAccountIdParam, strconv.Itoa(int(testAccount.ID)), 1), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestDeleteAccount_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			deleteAccount: func(accountId uint, dto *dto.DeleteAccountDto) error {
				return nil
			},
		},
	}
	router.DELETE(config.APIAccountIdPath, func(c echo.Context) error {
		login(container, c, testAccount)
		return account.DeleteAccount(c)
	})

	dto := dto.DeleteAccountDto{
		Password: testAccount.Password,
	}
	req := testutil.NewJSONRequest(http.MethodDelete, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID), dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteAccount_NoAuthorizationFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			deleteAccount: func(accountId uint, dto *dto.DeleteAccountDto) error {
				return nil
			},
		},
	}
	router.DELETE(config.APIAccountIdPath, func(c echo.Context) error {
		login(container, c, testAccount)
		return account.DeleteAccount(c)
	})

	dto := dto.DeleteAccountDto{
		Password: testAccount.Password,
	}
	req := testutil.NewJSONRequest(http.MethodDelete, fmt.Sprintf("%s/%d", config.APIAccount, testAccount.ID+1), dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	bodyBool, _ := strconv.ParseBool(rec.Body.String())
	assert.False(t, bodyBool)
}

func TestFindLoginId_Success(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			findAccountByEmail: func(dto *dto.FindLoginIdDto) error {
				return nil
			},
		},
	}
	router.POST(config.APIAccountFindLoginId, func(c echo.Context) error {
		return account.FindLoginId(c)
	})

	dto := dto.FindLoginIdDto{
		Email: testAccount.Email,
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccountFindLoginId, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFindLoginId_NoExistAccountFailure(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	testAccount := newTestUserAccount()
	account := accountController{
		container,
		&mockService{
			findAccountByEmail: func(dto *dto.FindLoginIdDto) error {
				return fmt.Errorf("account not found")
			},
		},
	}
	router.POST(config.APIAccountFindLoginId, func(c echo.Context) error {
		return account.FindLoginId(c)
	})

	dto := dto.FindLoginIdDto{
		Email: testAccount.Email,
	}
	req := testutil.NewJSONRequest(http.MethodPost, config.APIAccountFindLoginId, dto)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func login(testcontainer container.Container, c echo.Context, account model.Account) {
	_ = testcontainer.GetSession().Login(c,
		&infrastructure.Account{
			Id:        account.ID,
			LoginId:   account.LoginId,
			Authority: uint(account.Authority),
		},
	)
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
