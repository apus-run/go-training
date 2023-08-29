package dao

import (
	"context"
	"project-layout/internal/infra"
	"project-layout/internal/repository/dao/model"
)

// TopicDAO ... 数据访问层相关接口定义, 使用 DB 风格的命名
type TopicDAO interface {
	Insert(ctx context.Context, userModel model.Topic) (uint64, error)
	Update(ctx context.Context, userModel model.Topic) error
	Delete(ctx context.Context, userModel model.Topic) error
	FindByID(ctx context.Context, id uint64) (*model.Topic, error)
}

type topicDAO struct {
	data *infra.Data
}

func (t topicDAO) Insert(ctx context.Context, userModel model.Topic) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (t topicDAO) Update(ctx context.Context, userModel model.Topic) error {
	//TODO implement me
	panic("implement me")
}

func (t topicDAO) Delete(ctx context.Context, userModel model.Topic) error {
	//TODO implement me
	panic("implement me")
}

func (t topicDAO) FindByID(ctx context.Context, id uint64) (*model.Topic, error) {
	//TODO implement me
	panic("implement me")
}

func NewTopicDAO(data *infra.Data) TopicDAO {
	return &topicDAO{
		data: data,
	}
}
