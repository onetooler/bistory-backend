package dto

import "encoding/json"

type CreateAccountDto struct {
	LoginId  string `json:"loginId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewCreateAccountDto() *CreateAccountDto {
	return &CreateAccountDto{}
}

func (l *CreateAccountDto) ToString() (string, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

type LoginDto struct {
	LoginId  string `json:"loginId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginDto() *LoginDto {
	return &LoginDto{}
}

func (l *LoginDto) ToString() (string, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

type ChangeAccountPasswordDto struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"NewPassword"`
}

func NewChangeAccountPasswordDto() *ChangeAccountPasswordDto {
	return &ChangeAccountPasswordDto{}
}

func (l *ChangeAccountPasswordDto) ToString() (string, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

type DeleteAccountDto struct {
	Password string `json:"password"`
}

func NewDeleteAccountDto() *DeleteAccountDto {
	return &DeleteAccountDto{}
}

type FindLoginIdDto struct {
	Email string `json:"email"`
}

func NewFindLoginIdDto() *FindLoginIdDto {
	return &FindLoginIdDto{}
}
