package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 实体的属性对外部是不可见的
type User struct {
	ID       uint64
	Name     string // 账户名
	Email    string
	Password string
	Phone    string

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间

	ChangeTracker
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用
// ------------------------------------------------------------------------

func (u *User) setID(id uint64) *User {
	u.change()
	u.ID = id
	return u
}

func (u *User) setName(name string) *User {
	u.change()
	u.Name = name
	return u
}

func (u *User) setEmail(email string) *User {
	u.change()
	u.Email = email
	return u
}

func (u *User) setPassword(password string) *User {
	u.change()
	u.Password = password
	return u
}

func (u *User) setPhone(phone string) *User {
	u.change()
	u.Phone = phone
	return u
}

func (u *User) setCreatedTime(createdTime time.Time) *User {
	u.change()
	u.CreatedTime = createdTime
	return u
}

// 实体 JSON 序列化和反序列化
// ------------------------------------------------------------------------

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, u)
}

func (u *User) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

func (u *User) Scan(input any) error {
	return json.Unmarshal(input.([]byte), u)
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
