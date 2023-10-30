package service

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
)

// AuthService is a service for authentication.
type AuthService interface {
	AuthenticateByLoginIdAndPassword(loginId string, password string) (bool, *model.Account)
}

type authService struct {
	container container.Container
}

// NewAuthService is constructor.
func NewAuthService(container container.Container) AuthService {
	return &authService{container: container}
}

// AuthenticateByLoginIdAndPassword authenticates by using loginId and plain text password.
func (a *authService) AuthenticateByLoginIdAndPassword(loginId string, password string) (bool, *model.Account) {
	account, err := a.findByLoginId(loginId)
	if err != nil {
		a.container.GetLogger().GetZapLogger().Errorf(err.Error())
		return false, nil
	}

	if !account.CheckPassword(password) {
		a.container.GetLogger().GetZapLogger().Errorf("password not matched")
		return false, nil
	}

	return true, account
}

func (a *authService) findByLoginId(loginId string) (*model.Account, error) {
	repo := a.container.GetRepository()

	account := model.Account{LoginId: loginId}
	tx := repo.First(&account)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &account, nil
}
