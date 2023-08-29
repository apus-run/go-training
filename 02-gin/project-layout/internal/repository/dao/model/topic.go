package model

import (
	"database/sql/driver"
	"encoding/json"
	"project-layout/internal/domain/entity"
)

type Topic struct {
	ID           uint64 `gorm:"primaryKey,autoIncrement"`
	UserID       uint64 `gorm:"uniqueIndex;not null;comment:'用户ID'"`
	Title        string `gorm:"not null;comment:'主题标题'"`
	Content      string `gorm:"type:varchar(200) not null;comment:'主题内容'"`
	CommentCount int64  `gorm:"not null;comment:'评论数'"`

	Comments []*Comment
}

func (t *Topic) TableName() string {
	return "topic"
}

func (t *Topic) ToEntity() entity.Topic {
	if t == nil {
		return entity.Topic{}
	}
	// 构造comment
	comments := make([]entity.Comment, 0, len(t.Comments))

	for _, v := range t.Comments {
		commentBuilder := entity.NewCommentBuilder()
		commentBuilder.ID(v.ID)
		commentBuilder.TopicID(v.TopicID)
		commentBuilder.UserID(v.UserID)
		commentBuilder.Content(v.Content)
		commentEntity := commentBuilder.Build()
		comments = append(comments, *commentEntity)
	}

	// 构造topic
	topicBuilder := entity.NewTopicBuilder()
	topicBuilder.ID(t.ID)
	topicBuilder.UserID(t.UserID)
	topicBuilder.Title(t.Title)
	topicBuilder.Content(t.Content)
	topicBuilder.CommentCount(t.CommentCount)
	topicBuilder.Comment(comments)
	topicEntity := topicBuilder.Build()

	return *topicEntity

}

func (t *Topic) FromEntity(topicEntity entity.Topic) Topic {
	if t == nil {
		return Topic{}
	}

	t.ID = topicEntity.ID()
	t.UserID = topicEntity.UserID()
	t.Title = topicEntity.Title()
	t.Content = topicEntity.Content()
	t.CommentCount = topicEntity.CommentCount()
	return *t
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
