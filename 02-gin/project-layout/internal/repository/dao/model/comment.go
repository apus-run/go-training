package model

import (
	"database/sql/driver"
	"encoding/json"

	"project-layout/internal/domain/entity"
)

type Comment struct {
	ID      uint64 `gorm:"primaryKey,autoIncrement"`
	UserID  uint64 `gorm:"not null;comment:'用户ID'"`
	TopicID uint64 `gorm:"not null;comment:'主题ID'"`
	Content string `gorm:"type:varchar(500) not null;comment:'评论内容'"`
}

func (c *Comment) TableName() string {
	return "topic"
}

func (c *Comment) ToEntity() entity.Comment {
	commentBuilder := entity.NewCommentBuilder()
	commentBuilder.ID(c.ID)
	commentBuilder.TopicID(c.TopicID)
	commentBuilder.UserID(c.UserID)
	commentBuilder.Content(c.Content)
	commentEntity := commentBuilder.Build()

	return *commentEntity
}

func (c *Comment) FromEntity(commentEntity entity.Comment) Comment {
	if c == nil {
		return Comment{}
	}
	c.ID = commentEntity.ID()
	c.UserID = commentEntity.UserID()
	c.TopicID = commentEntity.TopicID()
	c.Content = commentEntity.Content()
	return *c
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
