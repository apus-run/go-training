package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"project-layout/internal/domain/entity"
)

type Topic struct {
	ID           uint64 `gorm:"primaryKey,autoIncrement"`
	UserID       uint64 `gorm:"uniqueIndex;not null;comment:'用户ID'"`
	Title        string `gorm:"not null;comment:'主题标题'"`
	Content      string `gorm:"type:varchar(200) not null;comment:'主题内容'"`
	CommentCount int64  `gorm:"not null;comment:'评论数'"`

	Comments []*Comment

	CreatedTime int64 `gorm:"autoCreateTime;not null;comment:'创建时间'"`
	UpdatedTime int64 `gorm:"autoUpdateTime;not null;comment:'更新时间'"`
	DeletedTime int64 `gorm:"index;not null;comment:'删除时间'"`
}

func (t *Topic) TableName() string {
	return "topic"
}

func (t *Topic) ToEntity() entity.Topic {
	// 构造comment
	comments := make([]*entity.Comment, 0, len(t.Comments))
	for _, v := range t.Comments {
		commentEntity := &entity.Comment{}
		commentEntity.ID = v.ID
		commentEntity.TopicID = v.TopicID
		commentEntity.UserID = v.UserID
		commentEntity.Content = v.Content
		comments = append(comments, commentEntity)
	}

	// 构造topic
	return entity.Topic{
		ID:           t.ID,
		UserID:       t.UserID,
		Title:        t.Title,
		Content:      t.Content,
		CommentCount: t.CommentCount,
		Comments:     comments,

		CreatedTime: time.UnixMilli(t.CreatedTime),
	}

}

func (t *Topic) FromEntity(topicEntity entity.Topic) Topic {
	return Topic{
		ID:           topicEntity.ID,
		UserID:       topicEntity.UserID,
		Title:        topicEntity.Title,
		Content:      topicEntity.Content,
		CommentCount: topicEntity.CommentCount,
	}
}

// MarshalBinary ...
func (t *Topic) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalBinary ...
func (t *Topic) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, t)
}

// Value ...
func (t *Topic) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

// Scan ...
func (t *Topic) Scan(input any) error {
	return json.Unmarshal(input.([]byte), t)
}
