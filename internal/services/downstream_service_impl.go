package services

import (
	"api-gateway/internal/global"
	"api-gateway/internal/model"
	"api-gateway/pkg/service"
	"context"

	"gorm.io/gorm"
)

type DownstreamServiceImpl struct {
	baseService service.BaseService[*model.Downstream]
}

func NewDownstreamService() DownstreamServiceImpl {
	bs := service.NewBaseService(&model.Downstream{}, global.DB)
	return DownstreamServiceImpl{
		baseService: bs,
	}
}

func (as *DownstreamServiceImpl) Add(ctx context.Context, data *model.Downstream) error {
	return as.baseService.Create(ctx, data)
}

func (as *DownstreamServiceImpl) GetByName(ctx context.Context, name string) (*model.Downstream, error) {
	return as.baseService.GetByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *DownstreamServiceImpl) GetByCondition(ctx context.Context, conditions map[string]any) ([]*model.Downstream, error) {
	return as.baseService.GetAllByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			tx = tx.Where(key, value)
		}
		return tx
	})
}

func (as *DownstreamServiceImpl) GetAll(ctx context.Context) ([]*model.Downstream, error) {
	return as.baseService.GetAllByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx
	})
}

func (as *DownstreamServiceImpl) Update(ctx context.Context, data model.Downstream) error {
	return as.baseService.UpdateById(ctx, &data)
}

func (as *DownstreamServiceImpl) UpdateByName(ctx context.Context, data model.Downstream, name string) error {
	return as.baseService.UpdateByCondition(ctx, &data, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *DownstreamServiceImpl) DeleteByName(ctx context.Context, name string) error {
	return as.baseService.DeleteByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *DownstreamServiceImpl) GetById(ctx context.Context, id uint) (*model.Downstream, error) {
	return as.baseService.GetByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}

func (as *DownstreamServiceImpl) Adds(ctx context.Context, datas []*model.Downstream) error {
	return as.baseService.CreateBatch(ctx, datas)
}

func (as *DownstreamServiceImpl) UpdateById(ctx context.Context, data model.Downstream, id uint) error {
	return as.baseService.UpdateByCondition(ctx, &data, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}

func (as *DownstreamServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return as.baseService.DeleteByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}
