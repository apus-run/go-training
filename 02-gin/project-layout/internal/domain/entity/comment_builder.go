package entity

type CommentBuilder struct {
	id      uint64
	topicID uint64
	userID  uint64
	content string
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{}
}

func (b *CommentBuilder) ID(id uint64) *CommentBuilder {
	b.id = id
	return b
}

func (b *CommentBuilder) TopicID(topicID uint64) *CommentBuilder {
	b.topicID = topicID
	return b
}

func (b *CommentBuilder) UserID(userID uint64) *CommentBuilder {
	b.userID = userID
	return b
}

func (b *CommentBuilder) Content(content string) *CommentBuilder {
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
