package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

const (
	MinLenPassword = 6
	MinLenLogin    = 1
	MaxLenLogin    = 25
)

//nolint:gochecknoinits
func init() {
	govalidator.CustomTypeTagMap.Set("login", func(i interface{}, o interface{}) bool {
		login, ok := i.(string)
		if !ok {
			return false
		}
		return len(login) >= MinLenLogin && len(login) <= MaxLenLogin
	})

	govalidator.CustomTypeTagMap.Set(
		"password",
		func(i interface{}, o interface{}) bool {
			pass, ok := i.(string)
			if !ok {
				return false
			}
			if len(pass) < MinLenPassword {
				return false
			}

			return true
		},
	)
}

type User struct {
	ID       uint64 `json:"id"        valid:"required"`
	Login    string `json:"login"     valid:"required,login"`
	Password string `json:"password"  valid:"required,password"`
	IsAdmin  bool   `json:"is_admin"  valid:"required"`
}

type UserWithoutPassword struct {
	ID      uint64 `json:"id"          valid:"required"`
	Login   string `json:"login"       valid:"required,login"`
	IsAdmin bool   `json:"is_admin"    valid:"required"`
}

func (u *UserWithoutPassword) Trim() {
	u.Login = strings.TrimSpace(u.Login)
}

type PreUser struct {
	Login    string `json:"login"    valid:"required,login"`
	Password string `json:"password" valid:"required,password"`
}

func (u *PreUser) Trim() {
	u.Login = strings.TrimSpace(u.Login)
}

func (u *UserWithoutPassword) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	u.Login = sanitizer.Sanitize(u.Login)
}
