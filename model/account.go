package model

import (
	"github.com/onetooler/bistory-backend/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Account defines struct of account data.
type Account struct {
	gorm.Model
	LoginId   string    `gorm:"unique;not null" json:"login_id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `json:"-"`
	Authority Authority `json:"authority"`
}

type Authority uint

const (
	AuthorityAdmin Authority = iota + 1
	AuthorityUser
)

func (a Authority) String() string {
	switch a {
	case AuthorityAdmin:
		return "Admin"
	case AuthorityUser:
		return "User"
	default:
		return "Invalid Authority"
	}
}

// TableName returns the table name of account struct and it is used by gorm.
func (Account) TableName() string {
	return "account"
}

func (a Account) CheckPassword(plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plainPassword)); err != nil {
		return false
	}
	return true
}

// NewAccountWithPasswordEncrypt is constructor. And it is encoded password by using bcrypt.
func NewAccountWithPasswordEncrypt(loginId, email, plainPassword string, authority Authority) (*Account, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), config.PasswordHashCost)
	if err != nil {
		return nil, err
	}
	return &Account{LoginId: loginId, Email: email, Password: string(hashed), Authority: authority}, nil
}

// ToString is return string of object
func (a *Account) ToString() string {
	return toString(a)
}
