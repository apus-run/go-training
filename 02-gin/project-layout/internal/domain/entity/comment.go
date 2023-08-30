package entity

type Comment struct {
	id      uint64
	topicID uint64
	userID  uint64
	content string

	ChangeTracker
}

// 实体的取值方法(get 关键字可以省略)
// 1、用于业务逻辑上需要取值的地方
// 2、用于基础设施层需要取值的地方
// ------------------------------------------------------------------------

func (c *Comment) ID() uint64 {
	return c.id
}

func (c *Comment) TopicID() uint64 {
	return c.topicID
}

func (c *Comment) UserID() uint64 {
	return c.userID
}

func (c *Comment) Content() string {
	return c.content
}

// 实体的赋值方法
// 实体的赋值方法可以变更属性，但是对外部也是不可见的，只能被实体的行为方法使用
// ------------------------------------------------------------------------

func (c *Comment) setID(id uint64) *Comment {
	c.change()
	c.id = id
	return c
}

func (c *Comment) setTopicID(topicID uint64) *Comment {
	c.change()
	c.topicID = topicID
	return c
}

func (c *Comment) setUserID(userID uint64) *Comment {
	c.change()
	c.userID = userID
	return c
}

func (c *Comment) setContent(content string) *Comment {
	c.change()
	c.content = content
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
