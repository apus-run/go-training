package entity

import "time"

type Comment struct {
	ID      uint64
	TopicID uint64
	UserID  uint64
	Content string

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间

	ChangeTracker
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用
// ------------------------------------------------------------------------

func (c *Comment) setID(id uint64) *Comment {
	c.change()
	c.ID = id
	return c
}

func (c *Comment) setTopicID(topicID uint64) *Comment {
	c.change()
	c.TopicID = topicID
	return c
}

func (c *Comment) setUserID(userID uint64) *Comment {
	c.change()
	c.UserID = userID
	return c
}

func (c *Comment) setContent(content string) *Comment {
	c.change()
	c.Content = content
	return c
}

// 实体行为方法
// 行为方法对外部是可见的，外部只能通过调用实体的行为方法来改变实体的属性
// ------------------------------------------------------------------------

// UpdateContent 变更评论内容
func (c *Comment) UpdateContent(content string) {
	// 判断字数不能超过多少字
	c.setContent(content)
}
