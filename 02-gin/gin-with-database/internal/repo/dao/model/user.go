package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"gin-with-database/internal/domain/entity"
)

// User ... sql.NullXXX 类型是为了方便处理 null 值; 空字符串是导致索引失效的罪魁祸首, 所以我们使用 sql.NullXXX
type User struct {
	ID uint64 `gorm:"primaryKey,autoIncrement"`

	Name string `gorm:"type:varchar(20) not null;comment:'用户名'"`

	// 设置邮箱为唯一索引
	Email    sql.NullString `gorm:"type:varchar(50);unique;not null;comment:'邮箱'"`
	Password string         `gorm:"type:varchar(150);not null;comment:'密码'"`
	// 设置手机号为唯一索引
	Phone sql.NullString `gorm:"type:varchar(20);unique;not null;comment:'手机号'"`

	CreatedTime int64 `gorm:"autoCreateTime;not null;comment:'创建时间'"`
	UpdatedTime int64 `gorm:"autoUpdateTime;not null;comment:'更新时间'"`
	DeletedTime int64 `gorm:"index;not null;comment:'删除时间'"`
}

// The TableName method returns the name of the data table that the struct is mapped to.
func (u *User) TableName() string {
	return "user"
}

func (u *User) ToEntity() entity.User {
	return entity.User{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email.String,
		Password:    u.Password,
		Phone:       u.Phone.String,
		CreatedTime: time.UnixMilli(u.CreatedTime),
	}
}

func (u *User) FromEntity(userEntity entity.User) User {
	return User{
		ID:   userEntity.ID,
		Name: userEntity.Name,
		Email: sql.NullString{
			String: userEntity.Email,
			Valid:  userEntity.Email != "",
		},
		Password: userEntity.Password,
		Phone: sql.NullString{
			String: userEntity.Phone,
			Valid:  userEntity.Phone != "",
		},
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
