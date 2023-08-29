package entity

type TopicBuilder struct {
	id           uint64
	userID       uint64
	title        string
	content      string
	commentCount int64
	comments     []Comment
}

func NewTopicBuilder() *TopicBuilder {
	return &TopicBuilder{}
}

func (b *TopicBuilder) ID(id uint64) *TopicBuilder {
	b.id = id
	return b
}

func (b *TopicBuilder) UserID(userID uint64) *TopicBuilder {
	b.userID = userID
	return b
}

func (b *TopicBuilder) Title(title string) *TopicBuilder {
	b.title = title
	return b
}

func (b *TopicBuilder) Content(content string) *TopicBuilder {
	b.content = content
	return b
}

func (b *TopicBuilder) CommentCount(commentCount int64) *TopicBuilder {
	b.commentCount = commentCount
	return b
}

func (b *TopicBuilder) Comment(comments []Comment) *TopicBuilder {
	b.comments = comments
	return b
}

func (b *TopicBuilder) Build() *Topic {
	return &Topic{
		id:           b.id,
		userID:       b.userID,
		title:        b.title,
		content:      b.content,
		commentCount: b.commentCount,
		comments:     b.comments,
	}
}
