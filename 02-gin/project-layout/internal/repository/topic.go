package repository

import (
	"context"
	"project-layout/internal/domain/entity"
	"project-layout/internal/repository/dao"
	"project-layout/pkg/log"
)

type TopicRepository interface {
	Save(ctx context.Context, entity entity.Topic) error

	Remove(ctx context.Context, entity entity.Topic) error
	Find(ctx context.Context, id int64) (*entity.Topic, error)

	FindByUserID(ctx context.Context, userID int64, lastTopicID int64) ([]*entity.Topic, error)
	FindCommentPageByUserId(ctx context.Context, topicID int64, userID int64, page int64, size int64) (*entity.Topic, int64, bool, error)
}

// topicRepository 使用了缓存
type topicRepository struct {
	dao dao.UserDAO

	log *log.Logger
}

func NewTopicRepository(dao dao.UserDAO, logger *log.Logger) TopicRepository {
	return &topicRepository{
		dao: dao,
		log: logger,
	}
}

func (tr *topicRepository) Save(ctx context.Context, topicEntity entity.Topic) (err error) {
	if topicEntity.ID() > 0 {
		err = tr.updateTopic(ctx, topicEntity)
		for _, commentEntity := range topicEntity.Comments() {
			if commentEntity.ID() > 0 {
				_ = tr.updateComment(ctx, commentEntity)
			} else {
				_ = tr.insertComment(ctx, commentEntity)
			}
		}
	} else {
		err = tr.insertTopic(ctx, topicEntity)
	}
	return nil
}

func (tr *topicRepository) insertTopic(ctx context.Context, topicEntity entity.Topic) (err error) {
	//insert topic 逻辑省略

	return nil
}

// 在 Infrastrcture 判断被标记有变动的才 Update
func (tr *topicRepository) updateTopic(ctx context.Context, topicEntity entity.Topic) (err error) {
	/*
		// Map the data from Entity to DO
		topicModel := new(model.Topic)
		err := topicModel.FromEntity(topicEntity)
		if err != nil {
			return err
		}
		// 在 Infrastrcture 判断被标记有变动的才 Update
		if !topicEntity.IsDirty() {
			return tr.data.DB.Transaction(func(tx *gorm.DB) error {
				// Create new record in the store
				err = tx.WithContext(ctx).Model(&topicModel).Updates(&topicModel).Error
				if err != nil {
					return err
				}
				// Map fresh record's data into Entity
				newEntity, err := topicModel.ToEntity()
				if err != nil {
					return err
				}
				*topicEntity = *newEntity
				return nil
			})
		}
	*/

	return nil
}

func (tr *topicRepository) insertComment(ctx context.Context, commentEntity entity.Comment) (err error) {
	//insert comment 逻辑省略

	return nil
}

// 在 Infrastrcture 判断被标记有变动的才 Update
func (tr *topicRepository) updateComment(ctx context.Context, commentEntity entity.Comment) (err error) {
	if !commentEntity.IsDirty() {
		return nil
	}

	//update comment 逻辑省略

	return nil
}

func (tr *topicRepository) Remove(ctx context.Context, entity entity.Topic) error {
	//TODO implement me
	panic("implement me")
}

func (tr *topicRepository) Find(ctx context.Context, id int64) (*entity.Topic, error) {
	//TODO implement me
	panic("implement me")
}

func (tr *topicRepository) FindByUserID(ctx context.Context, userID int64, lastTopicID int64) ([]*entity.Topic, error) {
	//TODO implement me
	panic("implement me")
}

func (tr *topicRepository) FindCommentPageByUserId(ctx context.Context, topicID int64, userID int64, page int64, size int64) (*entity.Topic, int64, bool, error) {
	//TODO implement me
	panic("implement me")
}
