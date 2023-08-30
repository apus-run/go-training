package entity

import (
	"project-layout/pkg/log"
)

type Topic struct {
	id           uint64
	userID       uint64
	title        string
	content      string
	commentCount int64
	comments     []*Comment

	ChangeTracker
}

// 实体的取值方法(get 关键字可以省略)
// 1、用于业务逻辑上需要取值的地方
// 2、用于基础设施层需要取值的地方
// ------------------------------------------------------------------------

func (t *Topic) ID() uint64 {
	return t.id
}

func (t *Topic) UserID() uint64 {
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
// ------------------------------------------------------------------------

func (t *Topic) setID(id uint64) *Topic {
	t.change()
	t.id = id
	return t
}

func (t *Topic) setUserID(userID uint64) *Topic {
	t.change()
	t.userID = userID
	return t
}

func (t *Topic) setContent(content string) *Topic {
	t.change()
	t.content = content
	return t
}

func (t *Topic) setCommentCount(commentCount int64) *Topic {
	t.change()
	t.commentCount = commentCount
	return t
}

func (t *Topic) setComments(comments []*Comment) *Topic {
	t.change()
	t.comments = comments
	return t
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性
// ------------------------------------------------------------------------

// AppendComment 聚合内一致性事务逻辑代码的实现
func (t *Topic) AppendComment(comment Comment) {
	t.setComments(append(t.Comments(), &comment))
	t.setCommentCount(t.CommentCount() + 1)
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
	t.setCommentCount(t.CommentCount() + num)
}

// DecreaseCommentCount 减少评论总数
func (t *Topic) DecreaseCommentCount(num int64) {
	// 判断字数不能少于多少
	if t.CommentCount()-num < 0 {
		log.Error("commentCount 不能超过最小值：0")
	}

	t.setCommentCount(t.CommentCount() - num)
}
