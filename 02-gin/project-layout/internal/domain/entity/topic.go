package entity

import (
	"project-layout/pkg/log"
	"time"
)

type Topic struct {
	ID           uint64
	UserID       uint64
	Title        string
	Content      string
	CommentCount int64
	Comments     []*Comment

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间

	ChangeTracker
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用
// ------------------------------------------------------------------------

func (t *Topic) setID(id uint64) *Topic {
	t.change()
	t.ID = id
	return t
}

func (t *Topic) setUserID(userID uint64) *Topic {
	t.change()
	t.UserID = userID
	return t
}

func (t *Topic) setContent(content string) *Topic {
	t.change()
	t.Content = content
	return t
}

func (t *Topic) setCommentCount(commentCount int64) *Topic {
	t.change()
	t.CommentCount = commentCount
	return t
}

func (t *Topic) setComments(comments []*Comment) *Topic {
	t.change()
	t.Comments = comments
	return t
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性
// ------------------------------------------------------------------------

// AppendComment 聚合内一致性事务逻辑代码的实现
func (t *Topic) AppendComment(comment Comment) {
	t.setComments(append(t.Comments, &comment))
	t.setCommentCount(t.CommentCount + 1)
}

// UpdateContent 变更话题内容
func (t *Topic) UpdateContent(content string) {
	// 判断字数不能超过多少字
	t.setContent(content)
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性

// UpdateCommentCount 变更话题下的评论总数
func (t *Topic) UpdateCommentCount(commentCount int64) {
	// 判断字数不能超过多少字
	t.setCommentCount(commentCount)
}

// IncreaseCommentCount 增加评论总数
func (t *Topic) IncreaseCommentCount(num int64) {
	t.setCommentCount(t.CommentCount + num)
}

// DecreaseCommentCount 减少评论总数
func (t *Topic) DecreaseCommentCount(num int64) {
	// 判断字数不能少于多少
	if t.CommentCount-num < 0 {
		log.Error("commentCount 不能超过最小值：0")
	}

	t.setCommentCount(t.CommentCount - num)
}
