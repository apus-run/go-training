package entity

type TopicBuilder struct {
	id           uint
	userID       string
	title        string
	content      string
	commentCount int64
	comments     []*Comment
}

func NewTopicBuilder() *TopicBuilder {
	return &TopicBuilder{}
}

func (b *TopicBuilder) WithID(id uint) *TopicBuilder {
	b.id = id
	return b
}

func (b *TopicBuilder) WithUserID(userID string) *TopicBuilder {
	b.userID = userID
	return b
}

func (b *TopicBuilder) WithTitle(title string) *TopicBuilder {
	b.title = title
	return b
}

func (b *TopicBuilder) WithContent(content string) *TopicBuilder {
	b.content = content
	return b
}

func (b *TopicBuilder) WithCommentCount(commentCount int64) *TopicBuilder {
	b.commentCount = commentCount
	return b
}

func (b *TopicBuilder) WithComment(comments []*Comment) *TopicBuilder {
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
