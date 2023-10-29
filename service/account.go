package service

import (
	"fmt"

	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
)

// AccountService is a service for managing user account.
type AccountService interface {
	CreateAccount(createAccountDto *dto.CreateAccountDto) (*model.Account, error)
	UpdateAccountPassword(loginId string, UpdatePasswordDto *dto.UpdatePasswordDto) (bool, error)
	DeleteAccount(loginId string) (bool, error)
	FindByLoginId(loginId string) (*model.Account, error)
}

type accountService struct {
	container container.Container
}

// NewAccountService is constructor.
func NewAccountService(container container.Container) AccountService {
	return &accountService{container: container}
}

func (a *accountService) CreateAccount(createAccountDto *dto.CreateAccountDto) (*model.Account, error) {
	// loginId validation
	exists, err := a.existsByLoginId(createAccountDto.LoginId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("loginId %s already exists", createAccountDto.LoginId)
	}

	// password validation
	if len(createAccountDto.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}
	if len(createAccountDto.Password) > 72 {
		return nil, fmt.Errorf("password must be at most 72 characters")
	}

	// create account
	account, err := model.NewAccountWithPasswordEncrypt(createAccountDto.LoginId, createAccountDto.Email, createAccountDto.Password, model.AuthorityUser)
	if err != nil {
		return nil, fmt.Errorf("create account failed: %s", err.Error())
	}

	err = a.create(account)
	if err != nil {
		return nil, err
	}
	
	return account, nil
}

func (a *accountService) UpdateAccountPassword(loginId string, UpdatePasswordDto *dto.UpdatePasswordDto) (bool, error) {
	// password validation
	if len(UpdatePasswordDto.Password) < 8 {
		a.container.GetLogger().GetZapLogger().Errorf("password must be at least 8 characters")
		return false, nil
	}
	if len(UpdatePasswordDto.Password) > 72 {
		a.container.GetLogger().GetZapLogger().Errorf("password must be at most 72 characters")
		return false, nil
	}

	return a.updatePassword(loginId, UpdatePasswordDto.Password)
}

func (a *accountService) DeleteAccount(loginId string) (bool, error) {
	repo := a.container.GetRepository()
	account := model.Account{LoginId:loginId}

	if err:=repo.Delete(&account).Error; err != nil {
		return false, err
	}
	return true, nil
}

// TODO: Need to review whether to change to ORM style call.
const existsAccount = "SELECT EXISTS (SELECT 1 FROM account WHERE id = ?);"

func (a *accountService) existsByLoginId(loginId string) (bool, error) {
	repo := a.container.GetRepository()

	exists := false
	tx := repo.Raw(existsAccount, loginId).Scan(&exists)

	return exists, tx.Error
}

func (a *accountService) FindByLoginId(loginId string) (*model.Account, error) {
	repo := a.container.GetRepository()

	account := model.Account{LoginId:loginId}
	tx := repo.First(&account)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &account, nil
}

func (a *accountService) create(account *model.Account) (error) {
	repo := a.container.GetRepository()

	if err := repo.Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (a *accountService) updatePassword(loginId string, password string) (bool, error) {
	repo := a.container.GetRepository()

	account := model.Account{LoginId:loginId}
	tx := repo.Model(account).Update("password", password)
	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil
}

