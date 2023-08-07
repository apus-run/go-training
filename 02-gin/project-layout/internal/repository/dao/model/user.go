package model

import (
	"database/sql/driver"
	"encoding/json"
	"project-layout/internal/domain/entity"
)

type User struct {
	ID uint64 `gorm:"primaryKey, autoIncrement"`

	Name     string
	Avatar   string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Phone    string

	CreatedTime uint64 // 创建时间
	UpdatedTime uint64 // 更新时间
	DeletedTime uint64 `gorm:"index"` // 删除时间
}

// The TableName method returns the name of the data table that the struct is mapped to.
func (u *User) TableName() string {
	return "user"
}

// ToEntity converts the DO to a entity.
func (u *User) ToEntity() (*entity.User, error) {
	user := &entity.User{}
	if u == nil {
		return user, nil
	}

	user.ID = u.ID
	user.Name = u.Name
	user.Email = u.Email
	user.Password = u.Password

	return user, nil
}

// FromEntity converts a entity to a DO.
func (u *User) FromEntity(userEntity *entity.User) error {
	if u == nil {
		u = &User{}
	}
	if err := userEntity.Validate(); err != nil {
		return err
	}

	u.ID = userEntity.ID
	u.Name = userEntity.Name
	u.Email = userEntity.Email
	u.Password = userEntity.Password

	return nil
}

// MarshalBinary ..
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// Value ..
func (u *User) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

// Scan ..
func (u *User) Scan(input any) error {
	return json.Unmarshal(input.([]byte), u)
}
