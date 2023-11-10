package service

import (
	"fmt"

	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
)

// AuthService is a service for authentication.
type AuthService interface {
	AuthenticateByLoginIdAndPassword(loginId string, password string) (*model.Account, error)
}

type authService struct {
	container container.Container
}

// NewAuthService is constructor.
func NewAuthService(container container.Container) AuthService {
	return &authService{container: container}
}

// AuthenticateByLoginIdAndPassword authenticates by using loginId and plain text password.
func (a *authService) AuthenticateByLoginIdAndPassword(loginId string, password string) (*model.Account, error) {
	account, err := a.findByLoginId(loginId)
	if err != nil {
		return nil, err
	}
	if !account.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}

	ok := account.CheckPassword(password)
	a.container.GetRepository().Save(account) // save
	if !ok {
		if account.RemainAttempt() > 0 {
			return nil, fmt.Errorf("password not matched. remain attempt count is %d", account.RemainAttempt())
		}
		return nil, fmt.Errorf("password not matched. account has been deactivated")
	}

	return account, nil
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
