package entity

import (
	"github.com/pkg/errors"
)

var ErrEmptyCommentID = errors.New("comment id is required")
var ErrEmptyCommentContent = errors.New("comment content is required")

type Comment struct {
	id      uint
	topicID uint
	userID  string
	content string
	ChangeTracker
}

func (c *Comment) ID() uint {
	return c.id
}

func (c *Comment) TopicID() uint {
	return c.topicID
}

func (c *Comment) UserID() string {
	return c.userID
}

func (c *Comment) Content() string {
	return c.content
}

func (c *Comment) withID(id uint) *Comment {
	c.id = id
	return c
}

func (c *Comment) withTopicID(topicID uint) *Comment {
	c.topicID = topicID
	return c
}

func (c *Comment) withUserID(userID string) *Comment {
	c.userID = userID
	return c
}

func (c *Comment) withContent(content string) *Comment {
	c.content = content
	return c
}

// Validate 参数校验
func (t *Comment) Validate() error {
	if t.topicID == 0 {
		return ErrEmptyTopicID
	}
	if t.userID == "" {
		return ErrEmptyUserID
	}

	if t.content == "" {
		return ErrEmptyCommentContent
	}

	return nil
}

func NewComment(topicID uint, userID, content string) (*Comment, error) {
	comment := &Comment{}
	return comment.
			withTopicID(topicID).
			withUserID(userID).
			withContent(content),
		nil
}
