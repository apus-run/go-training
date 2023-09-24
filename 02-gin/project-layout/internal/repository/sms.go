package repository

import (
	"context"
	"errors"
	"fmt"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository/dao"
	"project-layout/internal/repository/dao/model"
	"project-layout/pkg/log"
)

type SMSRepository interface {
	// Save 是 Create and Update 统称
	Save(ctx context.Context, smsEntity entity.Sms) error

	Remove(ctx context.Context, smsEntity entity.Sms) error
	FindByID(ctx context.Context, id uint64) (*entity.Sms, error)

	FindByStatus(ctx context.Context, status int32) ([]entity.Sms, error)
}

type smsRepository struct {
	dao dao.SmsDAO
	log *log.Logger
}

func NewSMSRepository(dao dao.SmsDAO, logger *log.Logger) SMSRepository {
	return &smsRepository{
		dao: dao,

		log: logger,
	}
}

// Save 保存用户信息到数据库
func (sr *smsRepository) Save(ctx context.Context, smsEntity entity.Sms) error {
	sr.log.Info("create or update sms")
	// Map the data from Entity to DO
	smsModel := model.Sms{}
	smsModel = smsModel.FromEntity(smsEntity)

	// 如果ID为0, 则是创建, 否则是更新
	if smsModel.ID == 0 {
		// Insert
		_, err := sr.dao.Insert(ctx, smsModel)
		if err != nil {
			return err
		}
	} else {
		// Update 更新数据，只有非 0 值才会更新
		err := sr.dao.Update(ctx, smsModel)
		if err != nil {
			return err
		}
	}

	// Map fresh record's data into Entity
	newEntity := smsModel.ToEntity()
	smsEntity = newEntity

	return nil
}

// Remove 从数据库和缓存中删除用户信息
func (sr *smsRepository) Remove(ctx context.Context, smsEntity entity.Sms) error {
	sr.log.Info("remove sms")
	// Map the data from Entity to DO
	smsModel := model.Sms{}
	smsModel = smsModel.FromEntity(smsEntity)

	// Remove the data from DB
	err := sr.dao.Delete(ctx, smsModel)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 从数据库和缓存中获取用户信息
func (sr *smsRepository) FindByID(ctx context.Context, id uint64) (*entity.Sms, error) {
	smsEntity := &entity.Sms{}

	userModel, err := sr.dao.FindByID(ctx, id)

	// 对错误进行包装, 尽量使用通一的错误处理(Sentinel error), 屏蔽掉不同数据库的报错差异性
	// 这里不只是返回这一种错误, 在上层打日志的时候还要看到底层出的是那种错误, 所以原始的错误信息也需要保留.
	if err != nil {
		if errors.Is(err, ErrUserDataNotFound) {
			return smsEntity, fmt.Errorf("记录不存在, err: %v", err)
		} else {
			return smsEntity, err
		}
	}

	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()

	return &newEntity, nil
}

func (sr *smsRepository) FindByStatus(ctx context.Context, status int32) ([]entity.Sms, error) {
	modelSmsList, err := sr.dao.FindByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	entitySmsList := make([]entity.Sms, 0, len(modelSmsList))
	for i, smsEntity := range modelSmsList {
		entitySmsList[i] = smsEntity.ToEntity()
	}

	return entitySmsList, nil
}
