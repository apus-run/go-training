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

// User 实体的属性对外部是不可见的
type User struct {
	id       uint64
	name     string // 账户名
	avatar   string
	email    string
	password string
	phone    string

	gender   int       // 性别
	nickName string    // 昵称
	realName string    // 真实姓名
	birthday time.Time // 生日
	profile  string    // 个人简介

	createdTime time.Time  // 创建时间
	updatedTime time.Time  // 更新时间
	deletedTime *time.Time // 删除时间

	ChangeTracker
}

// 实体的取值方法(get 关键字可以省略)

func (u *User) ID() uint64 {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Avatar() string {
	return u.avatar
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Phone() string {
	return u.phone
}

func (u *User) Gender() int {
	return u.gender
}

func (u *User) NickName() string {
	return u.nickName
}

func (u *User) RealName() string {
	return u.realName
}

func (u *User) Birthday() time.Time {
	return u.birthday
}

func (u *User) Profile() string {
	return u.profile
}

func (u *User) CreatedTime() time.Time {
	return u.createdTime
}

func (u *User) UpdatedTime() time.Time {
	return u.updatedTime
}

func (u *User) DeletedTime() *time.Time {
	return u.deletedTime
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用
// ------------------------------------------------------------------------

func (u *User) setID(id uint64) *User {
	u.change()
	u.id = id
	return u
}

func (u *User) setName(name string) *User {
	u.change()
	u.name = name
	return u
}

func (u *User) setAvatar(avatar string) *User {
	u.change()
	u.avatar = avatar
	return u
}

func (u *User) setEmail(email string) *User {
	u.change()
	u.email = email
	return u
}

func (u *User) setPassword(password string) *User {
	u.change()
	u.password = password
	return u
}

func (u *User) setPhone(phone string) *User {
	u.change()
	u.phone = phone
	return u
}

func (u *User) setGender(gender int) *User {
	u.change()
	u.gender = gender
	return u
}

func (u *User) setNickName(nickName string) *User {
	u.change()
	u.nickName = nickName
	return u
}

func (u *User) setRealName(realName string) *User {
	u.change()
	u.realName = realName
	return u
}

func (u *User) setBirthday(birthday time.Time) *User {
	u.change()
	u.birthday = birthday
	return u
}

func (u *User) setProfile(profile string) *User {
	u.change()
	u.profile = profile
	return u
}

func (u *User) setCreatedTime(createdTime time.Time) *User {
	u.change()
	u.createdTime = createdTime
	return u
}

func (u *User) Validate() error {
	if len(u.name) == 0 {
		return ErrInvalidUserName
	}

	if len(u.password) == 0 {
		return ErrInvalidPassword
	}

	ok, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, u.email)
	if !ok || err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性
// ------------------------------------------------------------------------

// UpdatePassword 变更密码
func (u *User) UpdatePassword(password string) {
	u.setPassword(password)
}

// UpdateEmail 变更邮箱
func (u *User) UpdateEmail(email string) {
	u.setEmail(email)
}

// UpdatePhone 变更手机号
func (u *User) UpdatePhone(phone string) {
	u.setPhone(phone)
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
