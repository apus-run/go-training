package dao

import (
	"context"
	"time"

	"project-layout/internal/infra"
	"project-layout/internal/repository/dao/model"
)

type SmsDAO interface {
	Insert(ctx context.Context, smsModel model.Sms) (uint64, error)
	Update(ctx context.Context, smsModel model.Sms) error
	Delete(ctx context.Context, smsModel model.Sms) error
	FindByID(ctx context.Context, id uint64) (*model.Sms, error)
	FindByStatus(ctx context.Context, status int32) ([]model.Sms, error)
}

type smsDAO struct {
	data *infra.Data
}

func NewSmsDAO(data *infra.Data) SmsDAO {
	return &smsDAO{
		data: data,
	}
}

func (s *smsDAO) Insert(ctx context.Context, smsModel model.Sms) (uint64, error) {
	now := time.Now().UnixMilli()
	smsModel.CreatedTime = now
	smsModel.UpdatedTime = now

	err := s.data.DB.WithContext(ctx).Create(&smsModel).Error

	return smsModel.ID, err
}

func (s *smsDAO) Update(ctx context.Context, smsModel model.Sms) error {
	// 这种写法是很不清晰的，因为它依赖了 gorm 的两个默认语义
	// 会使用 ID 来作为 WHERE 条件
	// 会使用非零值来更新
	// 另外一种做法是显式指定只更新必要的字段，
	// 那么这意味着 DAO 和 service 中非敏感字段语义耦合了
	err := s.data.DB.WithContext(ctx).Updates(&smsModel).Error
	return err
}

func (s *smsDAO) Delete(ctx context.Context, smsModel model.Sms) error {
	return s.data.DB.WithContext(ctx).Delete(&smsModel, "id = ?", smsModel.ID).Error
}

func (s *smsDAO) FindByID(ctx context.Context, id uint64) (*model.Sms, error) {
	var smsModel model.Sms
	err := s.data.DB.WithContext(ctx).First(&smsModel, "id = ?", id).Error
	return &smsModel, err
}

func (s *smsDAO) FindByStatus(ctx context.Context, status int32) ([]model.Sms, error) {
	result := make([]model.Sms, 0, 10)

	s.data.DB.WithContext(ctx).Find(&result, "status = ?", status)

	return result, nil
}
