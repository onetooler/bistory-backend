package service

import (
	"testing"

	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAccountCreate_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest()
	service := NewAccountService(container)

	createDto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTestTest",
	}
	account, err := service.CreateAccount(&createDto)
	assert.Nil(t, err)

	// auto-generated check
	assert.NotEmpty(t, account.ID)
	assert.NotEmpty(t, account.CreatedAt)
	assert.NotEmpty(t, account.UpdatedAt)
	assert.Empty(t, account.DeletedAt)
	assert.Equal(t, model.AuthorityUser, account.Authority)

	// equal check
	assert.Equal(t, createDto.LoginId, account.LoginId)
	assert.Equal(t, createDto.Email, account.Email)
	assert.True(t, account.CheckPassword(createDto.Password))
}

func TestAccountCreate_WrongPasswordFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest()
	service := NewAccountService(container)

	createDto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTest",
	}
	account, err := service.CreateAccount(&createDto)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestAccountGet_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest()
	service := NewAccountService(container)
	savedAccount := createSuccessAccount(service)

	account, err := service.GetAccount(savedAccount.ID)
	assert.Nil(t, err)

	account.CreatedAt = account.CreatedAt.Local()
	account.UpdatedAt = account.UpdatedAt.Local()
	assert.EqualValues(t, savedAccount, account)
}

func TestAccountGet_NotExistsIdFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest()
	service := NewAccountService(container)

	account, err := service.GetAccount(uint(999))
	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestChangeAccountPassword_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAccountService(container)
	savedAccount := createSuccessAccount(service)

	changeAccountPasswordDto := dto.ChangeAccountPasswordDto{
		OldPassword: "newTestTest",
		NewPassword: "newTestTestTest",
	}
	account, err := service.ChangeAccountPassword(savedAccount.ID, &changeAccountPasswordDto)
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.NotEqual(t, savedAccount.UpdatedAt, account.UpdatedAt)
}

func TestChangeAccountPassword_WrongPasswordFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAccountService(container)
	savedAccount := createSuccessAccount(service)

	changeAccountPasswordDto := dto.ChangeAccountPasswordDto{
		OldPassword: "newTestTest",
		NewPassword: "new",
	}
	account, err := service.ChangeAccountPassword(savedAccount.ID, &changeAccountPasswordDto)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestDeleteAccount_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAccountService(container)
	savedAccount := createSuccessAccount(service)

	dto := dto.DeleteAccountDto{
		Password: "newTestTest",
	}
	err := service.DeleteAccount(savedAccount.ID, &dto)
	assert.Nil(t, err)

	account, err := service.GetAccount(savedAccount.ID)
	assert.Nil(t, account)
	assert.NotNil(t, err)
}

func TestDeleteAccount_WrongPasswordFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAccountService(container)
	savedAccount := createSuccessAccount(service)

	dto := dto.DeleteAccountDto{
		Password: "newTest",
	}
	err := service.DeleteAccount(savedAccount.ID, &dto)
	assert.NotNil(t, err)
}

func createSuccessAccount(service AccountService) *model.Account {
	createDto := dto.CreateAccountDto{
		LoginId:  "newTest",
		Email:    "newTest@example.com",
		Password: "newTestTest",
	}
	savedAccount, _ := service.CreateAccount(&createDto)
	return savedAccount
}
