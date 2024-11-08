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

func (srv *DownstreamServiceImpl) AddService(ctx context.Context, data *model.Downstream) error {
	return srv.baseService.Create(ctx, data)
}

func (srv *DownstreamServiceImpl) GetService(ctx context.Context, name string) (model.Downstream, error) {
	var data model.Downstream
	err := srv.baseService.GetByCondition(ctx, &data, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})

	return data, err
}

func (srv *DownstreamServiceImpl) RemoveService(ctx context.Context, name string) error {
	return srv.baseService.DeleteByCondition(ctx, &model.Downstream{}, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}
