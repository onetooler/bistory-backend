package service

import (
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/onetooler/bistory-backend/util"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateByLoginIdAndPassword_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest(false)

	service := NewAuthService(container)
	account, err := service.AuthenticateByLoginIdAndPassword("test", "test")
	account.CreatedAt = account.CreatedAt.Local()
	account.UpdatedAt = account.UpdatedAt.Local()

	data := model.Account{LoginId: "test"}
	container.GetRepository().First(&data)
	data.CreatedAt = data.CreatedAt.Local()
	data.UpdatedAt = data.UpdatedAt.Local()

	assert.Equal(t, data, *account)
	assert.Nil(t, err)
}

func TestAuthenticateByLoginIdAndPassword_EntityNotFound(t *testing.T) {
	container := testutil.PrepareForServiceTest(false)

	service := NewAuthService(container)
	account, err := service.AuthenticateByLoginIdAndPassword("abcde", "abcde")

	assert.Nil(t, account)
	assert.NotNil(t, err)
}

func TestAuthenticateByLoginIdAndPassword_AuthenticationFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest(false)

	service := NewAuthService(container)
	account, err := service.AuthenticateByLoginIdAndPassword("test", "abcde")

	assert.Nil(t, account)
	assert.NotNil(t, err)
}

func TestAuthenticateByLoginIdAndPassword_AuthenticationMaxFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest(false)

	service := NewAuthService(container)
	for i := 0; i < config.MaxLoginAttempts; i++ {
		account, err := service.AuthenticateByLoginIdAndPassword("test", "abcde")
		assert.Nil(t, account)
		assert.NotNil(t, err)
	}
	account, err := service.AuthenticateByLoginIdAndPassword("test", "test")
	assert.Nil(t, account)
	assert.NotNil(t, err)
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

	container := testutil.PrepareForServiceTest(true)
	service := NewAuthService(container)

	testEmail := "testEmail@example.com"
	token, err := service.EmailVerificationTokenSend(testEmail)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Contains(t, mailServer.Messages()[1].MsgRequest(), *token)
}
