package service

import (
	"fmt"

	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/model/dto"
	"gorm.io/gorm/clause"
)

// AccountService is a service for managing user account.
type AccountService interface {
	CreateAccount(*dto.CreateAccountDto) (*model.Account, error)
	GetAccount(id uint) (*model.Account, error)
	ChangeAccountPassword(uint, *dto.ChangeAccountPasswordDto) (*model.Account, error)
	DeleteAccount(id uint) (bool, error)
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
	if err := a.validatePassword(createAccountDto.Password); err != nil {
		return nil, err
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

func (a *accountService) GetAccount(id uint) (*model.Account, error) {
	repo := a.container.GetRepository()

	account := model.Account{}
	tx := repo.First(&account, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &account, nil
}

func (a *accountService) ChangeAccountPassword(id uint, changeAccountPasswordDto *dto.ChangeAccountPasswordDto) (*model.Account, error) {
	// OldPassword validation
	account, err := a.GetAccount(id)
	if err != nil {
		return nil, err
	}
	ok := account.CheckPassword(changeAccountPasswordDto.OldPassword)
	if !ok {
		return nil, fmt.Errorf("old password is not valid")
	}

	// NewPassword validation
	if err := a.validatePassword(changeAccountPasswordDto.NewPassword); err != nil {
		return nil, err
	}

	return a.updatePassword(id, changeAccountPasswordDto.NewPassword)
}

func (a *accountService) DeleteAccount(id uint) (bool, error) {
	repo := a.container.GetRepository()

	if err := repo.Delete(&model.Account{}, id).Error; err != nil {
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

func (a *accountService) GetAccountByLoginId(loginId string) (*model.Account, error) {
	repo := a.container.GetRepository()

	account := model.Account{LoginId: loginId}
	tx := repo.First(&account)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &account, nil
}

func (a *accountService) create(account *model.Account) error {
	repo := a.container.GetRepository()

	if err := repo.Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (a *accountService) updatePassword(id uint, password string) (*model.Account, error) {
	repo := a.container.GetRepository()

	account := model.Account{}
	account.ID = id
	tx := repo.Model(&account).Clauses(clause.Returning{}).Update("password", password)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &account, nil
}

func (a *accountService) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters")
	}
	return nil
}
