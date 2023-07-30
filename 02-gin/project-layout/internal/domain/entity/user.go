package entity

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"
)

var ErrEmptyUserID = errors.New("user id is required")
var ErrEmptyUserName = errors.New("user name is required")
var ErrEmptyUserEmail = errors.New("user email is required")
var ErrEmptyUserPassword = errors.New("user password is required")
var ErrInvalidUserName = errors.New("用户名不合法")
var ErrInvalidPassword = errors.New("密码不合法")
var ErrInvalidEmail = errors.New("邮箱不合法")

type Users []User
type User struct {
	ID       uint64
	Name     string
	Avatar   string
	Email    string
	Password string
	Salt     string
}

func (u *User) Validate() error {
	if len(u.Name) == 0 {
		return ErrInvalidUserName
	}

	if len(u.Password) == 0 {
		return ErrInvalidPassword

	}

	ok, err := regexp.Match("", []byte(u.Email))
	if !ok || err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// HashPassword 密码加密
func (u *User) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	pass := string(b)
	return pass, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}