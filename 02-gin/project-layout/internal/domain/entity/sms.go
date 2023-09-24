package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Sms struct {
	ID uint64

	Biz     string // 业务签名
	Args    string // 参数
	Numbers string // 发送号码
	Status  int32  // 状态, 1: 发送成功, 2: 发送失败

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间

	ChangeTracker
}

func (s *Sms) setStatus(status int32) *Sms {
	s.change()
	s.Status = status
	return s
}

func (s *Sms) setCreatedTime(createdTime time.Time) *Sms {
	s.change()
	s.CreatedTime = createdTime
	return s
}

// 实体 JSON 序列化和反序列化
// ------------------------------------------------------------------------

func (s *Sms) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Sms) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s *Sms) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *Sms) Scan(input any) error {
	return json.Unmarshal(input.([]byte), s)
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性
// ------------------------------------------------------------------------

// UpdateStatus 更新状态
func (s *Sms) UpdateStatus(status int32) {
	s.setStatus(status)
}
