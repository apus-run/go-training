package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"project-layout/internal/domain/entity"
)

type Comment struct {
	ID      uint64 `gorm:"primaryKey,autoIncrement"`
	UserID  uint64 `gorm:"not null;comment:'用户ID'"`
	TopicID uint64 `gorm:"not null;comment:'主题ID'"`
	Content string `gorm:"type:varchar(500) not null;comment:'评论内容'"`

	CreatedTime int64 `gorm:"autoCreateTime;not null;comment:'创建时间'"`
	UpdatedTime int64 `gorm:"autoUpdateTime;not null;comment:'更新时间'"`
	DeletedTime int64 `gorm:"index;not null;comment:'删除时间'"`
}

func (c *Comment) TableName() string {
	return "topic"
}

func (c *Comment) ToEntity() entity.Comment {
	return entity.Comment{
		ID:      c.ID,
		TopicID: c.TopicID,
		UserID:  c.UserID,
		Content: c.Content,

		CreatedTime: time.UnixMilli(c.CreatedTime),
	}
}

func (c *Comment) FromEntity(commentEntity entity.Comment) Comment {
	return Comment{
		ID:      commentEntity.ID,
		UserID:  commentEntity.UserID,
		TopicID: commentEntity.TopicID,
		Content: commentEntity.Content,
	}
}

// MarshalBinary ...
func (c *Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// UnmarshalBinary ...
func (c *Comment) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

// Value ...
func (c *Comment) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

// Scan ...
func (c *Comment) Scan(input any) error {
	return json.Unmarshal(input.([]byte), c)
}
