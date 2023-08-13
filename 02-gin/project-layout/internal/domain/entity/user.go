package entity

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	Name     string // 账户名
	Avatar   string
	Email    string
	Password string
	Phone    string

	Gender   int    // 性别
	NickName string // 昵称
	RealName string // 真实姓名
	Birthday string // 生日
	Profile  string // 个人简介

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间
}

func (u *User) Validate() error {
	if len(u.Name) == 0 {
		return ErrInvalidUserName
	}

	if len(u.Password) == 0 {
		return ErrInvalidPassword

	}

	ok, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, u.Email)
	if !ok || err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// GenerateHashPassword 密码加密
func (u *User) GenerateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	pass := string(hash)
	return pass, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
