package entity

type CommentBuilder struct {
	id      uint
	topicID uint
	userID  string
	content string
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{}
}

func (b *CommentBuilder) WithID(id uint) *CommentBuilder {
	b.id = id
	return b
}

func (b *CommentBuilder) WithTopicID(topicID uint) *CommentBuilder {
	b.topicID = topicID
	return b
}

func (b *CommentBuilder) WithUserID(userID string) *CommentBuilder {
	b.userID = userID
	return b
}

func (b *CommentBuilder) WithContent(content string) *CommentBuilder {
	b.content = content
	return b
}

func (b *CommentBuilder) Build() *Comment {
	return &Comment{
		id:      b.id,
		topicID: b.topicID,
		userID:  b.userID,
		content: b.content,
	}
}
