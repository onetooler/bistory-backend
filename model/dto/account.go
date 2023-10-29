package dto

import "encoding/json"

type CreateAccountDto struct {
	LoginId string `json:"loginId"`
	Email string `json:"email"`
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
	LoginId string `json:"loginId"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewLoginDto() *LoginDto {
	return &LoginDto{}
}

func (l *LoginDto) ToString() (string, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}


type UpdatePasswordDto struct {
	Password string `json:"password"`
}

func NewUpdatePasswordDto() *UpdatePasswordDto {
	return &UpdatePasswordDto{}
}

func (l *UpdatePasswordDto) ToString() (string, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}
