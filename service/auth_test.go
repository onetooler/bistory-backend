package service

import (
	"testing"

	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateByLoginIdAndPassword_Success(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAuthService(container)
	result, account := service.AuthenticateByLoginIdAndPassword("test", "test")

	data := model.Account{LoginId: "test"}
	container.GetRepository().First(&data)

	assert.Equal(t, data, *account)
	assert.True(t, result)
}

func TestAuthenticateByLoginIdAndPassword_EntityNotFound(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAuthService(container)
	result, account := service.AuthenticateByLoginIdAndPassword("abcde", "abcde")

	assert.Nil(t, account)
	assert.False(t, result)
}

func TestAuthenticateByLoginIdAndPassword_AuthenticationFailure(t *testing.T) {
	container := testutil.PrepareForServiceTest()

	service := NewAuthService(container)
	result, account := service.AuthenticateByLoginIdAndPassword("test", "abcde")

	assert.Nil(t, account)
	assert.False(t, result)
}
