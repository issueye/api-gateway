package services

import (
	"api-gateway/pkg/service"
	"context"

	"api-gateway/internal/global"
	"api-gateway/internal/model"

	"gorm.io/gorm"
)

type APIServiceImpl struct {
	baseService service.BaseService[*model.APIInfo]
}

func NewAPIService() APIServiceImpl {
	bs := service.NewBaseService(&model.APIInfo{}, global.DB)
	return APIServiceImpl{
		baseService: bs,
	}
}

func (as *APIServiceImpl) Add(ctx context.Context, apiInfo *model.APIInfo) error {
	return as.baseService.Create(ctx, apiInfo)
}

func (as *APIServiceImpl) GetByName(ctx context.Context, name string) (*model.APIInfo, error) {
	return as.baseService.GetByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *APIServiceImpl) GetByCondition(ctx context.Context, conditions map[string]any) ([]*model.APIInfo, error) {
	return as.baseService.GetAllByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			tx = tx.Where(key, value)
		}
		return tx
	})
}

func (as *APIServiceImpl) GetAll(ctx context.Context) ([]*model.APIInfo, error) {
	return as.baseService.GetAllByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx
	})
}

func (as *APIServiceImpl) Update(ctx context.Context, apiInfo model.APIInfo) error {
	return as.baseService.UpdateById(ctx, &apiInfo)
}

func (as *APIServiceImpl) UpdateByName(ctx context.Context, apiInfo model.APIInfo, name string) error {
	return as.baseService.UpdateByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *APIServiceImpl) DeleteByName(ctx context.Context, name string) error {
	return as.baseService.DeleteByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *APIServiceImpl) GetById(ctx context.Context, id uint) (*model.APIInfo, error) {
	return as.baseService.GetByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}

func (as *APIServiceImpl) Adds(ctx context.Context, apiInfos []*model.APIInfo) error {
	return as.baseService.CreateBatch(ctx, apiInfos)
}

func (as *APIServiceImpl) UpdateById(ctx context.Context, apiInfo model.APIInfo, id uint) error {
	return as.baseService.UpdateByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}

func (as *APIServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return as.baseService.DeleteByCondition(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}
