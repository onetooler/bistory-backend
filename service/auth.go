package service

import (
	"fmt"

	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/util"
)

// AuthService is a service for authentication.
type AuthService interface {
	AuthenticateByLoginIdAndPassword(loginId string, password string) (*model.Account, error)
	EmailVerificationTokenSend(email string) (*string, error)
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

// EmailVerificationTokenSend send token to email and return that token.
func (a *authService) EmailVerificationTokenSend(email string) (*string, error) {
	emailSender := a.container.GetEmailSender()
	token := util.RandomBase16String(config.EmailVerificationTokenLength)
	// TODO: Change to Constant
	subject := "[Bistory] 이메일 인증 코드"
	err := emailSender.SendEmail(email, subject, config.EmailVerificationTemplate, token)
	if err != nil {
		return nil, err
	}
	return &token, nil
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
