package entity

import (
	"log"

	"github.com/pkg/errors"
)

var ErrEmptyTopicID = errors.New("topic id is required")
var ErrEmptyTopicTitle = errors.New("topic title is required")
var ErrEmptyTopicContent = errors.New("topic content is required")
var ErrEmptyTopicCommentCount = errors.New("topic comment count is required")
var ErrEmptyTopicComments = errors.New("topic comments is required")

type Topic struct {
	id           uint
	userID       string
	title        string
	content      string
	commentCount int64
	comments     []*Comment

	ChangeTracker
}

func (t *Topic) ID() uint {
	return t.id
}

func (t *Topic) UserID() string {
	return t.userID
}

func (t *Topic) Title() string {
	return t.title
}

func (t *Topic) Content() string {
	return t.content
}

func (t *Topic) CommentCount() int64 {
	return t.commentCount

}

func (t *Topic) Comments() []*Comment {
	return t.comments
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用

func (t *Topic) withID(id uint) *Topic {
	t.change()
	t.id = id
	return t
}

func (t *Topic) withUserID(userID string) *Topic {
	t.change()
	t.userID = userID
	return t
}

func (t *Topic) withContent(content string) *Topic {
	t.change()
	t.content = content
	return t
}

func (t *Topic) withCommentCount(commentCount int64) *Topic {
	t.change()
	t.commentCount = commentCount
	return t
}

func (t *Topic) withComments(comments []*Comment) *Topic {
	t.change()
	t.comments = comments
	return t
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性

// AppendComment 聚合内一致性事务逻辑代码的实现
func (t *Topic) AppendComment(comment *Comment) {
	t.withComments(append(t.Comments(), comment))
	t.withCommentCount(t.CommentCount() + 1)
}

// UpdateContent 变更话题内容
func (t *Topic) UpdateContent(content string) {
	// 判断字数不能超过多少字
	t.withContent(content)
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性

// UpdateCommentCount 变更话题下的评论总数
func (t *Topic) UpdateCommentCount(commentCount int64) {
	// 判断字数不能超过多少字
	t.withCommentCount(commentCount)
}

// IncreaseCommentCount 增加评论总数
func (t *Topic) IncreaseCommentCount(num int64) {
	t.withCommentCount(t.CommentCount() + num)
}

// DecreaseCommentCount 减少评论总数
func (t *Topic) DecreaseCommentCount(num int64) {
	// 判断字数不能少于多少
	if t.CommentCount()-num < 0 {
		log.Panicf("commentCount 不能超过最小值：0")
	}

	t.withCommentCount(t.CommentCount() - num)
}

// Validate 参数校验
func (t *Topic) Validate() error {
	if t.userID == "" {
		return ErrEmptyUserID
	}
	if t.title == "" {
		return ErrEmptyTopicTitle
	}

	if t.content == "" {
		return ErrEmptyTopicContent
	}
	if t.commentCount == 0 {
		return ErrEmptyTopicCommentCount
	}

	if len(t.comments) == 0 {
		return ErrEmptyTopicComments
	}

	return nil
}

func NewTopic(userID, content string, commentCount int64, comments []*Comment) (*Topic, error) {
	topic := &Topic{}
	err := topic.Validate()
	if err != nil {
		return topic, err
	}

	return topic.
		withUserID(userID).
		withContent(content).
		withCommentCount(commentCount).
		withComments(comments), nil
}
