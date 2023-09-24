package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"project-layout/internal/domain/entity"
)

type Sms struct {
	ID      uint64 `gorm:"primaryKey,autoIncrement"`
	Biz     string `gorm:"comment:'业务签名"`
	Args    string `gorm:"comment:'参数"`
	Numbers string `gorm:"comment:'发送的号码"`
	Status  int32  `gorm:"comment:'状态 1 发送成功 2 发送失败"`

	CreatedTime int64 `gorm:"autoCreateTime;not null;comment:'创建时间'"`
	UpdatedTime int64 `gorm:"autoUpdateTime;not null;comment:'更新时间'"`
	DeletedTime int64 `gorm:"index;not null;comment:'删除时间'"`
}

func (s *Sms) TableName() string {
	return "sms"
}

func (s *Sms) ToEntity() entity.Sms {
	return entity.Sms{
		ID:      s.ID,
		Biz:     s.Biz,
		Args:    s.Args,
		Numbers: s.Numbers,
		Status:  s.Status,

		CreatedTime: time.UnixMilli(s.CreatedTime),
	}
}

func (s *Sms) FromEntity(smsEntity entity.Sms) Sms {
	return Sms{
		ID:      smsEntity.ID,
		Biz:     smsEntity.Biz,
		Args:    smsEntity.Args,
		Numbers: smsEntity.Numbers,
		Status:  smsEntity.Status,
	}
}

// MarshalBinary ...
func (s *Sms) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary ...
func (s *Sms) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

// Value ...
func (s *Sms) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

// Scan ...
func (s *Sms) Scan(input any) error {
	return json.Unmarshal(input.([]byte), s)
}
