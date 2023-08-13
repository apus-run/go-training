package model

import (
	"database/sql/driver"
	"encoding/json"

	"project-layout/internal/domain/entity"
)

var _ Model[entity.User] = (*User)(nil)

type User struct {
	ID uint64 `gorm:"primaryKey,autoIncrement"`

	Name     string `gorm:"type:varchar(20) not null;display:'';comment:'用户名'"`
	Avatar   string `gorm:"type:varchar(200);not null;display:'';comment:'头像'"`
	Email    string `gorm:"type:varchar(50);uniqueIndex;not null;display:'';comment:'邮箱'"` // 设置邮箱为唯一索引
	Password string `gorm:"type:varchar(150);not null;display:'';comment:'密码'"`
	Phone    string `gorm:"type:varchar(20);unique;not null;display:'';comment:'手机号'"`

	Gender   int    `gorm:"display:1;comment:'性别'"`
	NickName string `gorm:"type:varchar(20);display:'';comment:'昵称'"`
	RealName string `gorm:"type:varchar(20);display:'';comment:'真实姓名'"`
	Birthday string `gorm:"type:varchar(20);display:'';comment:'生日'"`
	Profile  string `gorm:"type:varchar(255);display:'';comment:'个人简介'"`

	CreatedTime int64 `gorm:"autoCreateTime;not null;comment:'创建时间'"`
	UpdatedTime int64 `gorm:"autoUpdateTime;not null;comment:'更新时间'"`
	DeletedTime int64 `gorm:"index;not null;display:0;comment:'删除时间'"`
}

// The TableName method returns the name of the data table that the struct is mapped to.
func (u *User) TableName() string {
	return "user"
}

func (u *User) ToEntity() entity.User {
	if u == nil {
		return entity.User{}
	}
	return entity.User{
		ID:       u.ID,
		Name:     u.Name,
		Avatar:   u.Avatar,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,

		Gender:   u.Gender,
		NickName: u.NickName,
		RealName: u.RealName,
		Birthday: u.Birthday,
		Profile:  u.Profile,
	}
}

func (u *User) FromEntity(userEntity entity.User) any {
	if u == nil {
		return User{}
	}
	if err := userEntity.Validate(); err != nil {
		return err
	}
	return User{
		ID:       userEntity.ID,
		Name:     userEntity.Name,
		Avatar:   userEntity.Avatar,
		Email:    userEntity.Email,
		Password: userEntity.Password,
		Phone:    userEntity.Phone,

		Gender:   userEntity.Gender,
		NickName: userEntity.NickName,
		RealName: userEntity.RealName,
		Birthday: userEntity.Birthday,
		Profile:  userEntity.Profile,
	}
}

// MarshalBinary ...
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *User) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, u)
}

// Value ...
func (u *User) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

// Scan ...
func (u *User) Scan(input any) error {
	return json.Unmarshal(input.([]byte), u)
}
