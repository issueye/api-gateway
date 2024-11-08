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

func (as *APIServiceImpl) AddAPIInfo(ctx context.Context, apiInfo *model.APIInfo) error {
	return as.baseService.Create(ctx, apiInfo)
}

func (as *APIServiceImpl) GetAPIInfo(ctx context.Context, name string) (model.APIInfo, error) {
	var apiInfo model.APIInfo
	err := as.baseService.GetByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
	return apiInfo, err
}

func (as *APIServiceImpl) GetAPIInfos(ctx context.Context, conditions map[string]any) ([]*model.APIInfo, error) {
	apiInfos := make([]*model.APIInfo, 0)
	err := as.baseService.GetAllByCondition(ctx, apiInfos, func(tx *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			tx = tx.Where(key, value)
		}
		return tx
	})
	return apiInfos, err
}

func (as *APIServiceImpl) GetAllAPIInfo(ctx context.Context) ([]*model.APIInfo, error) {
	apiInfos := make([]*model.APIInfo, 0)
	err := as.baseService.GetAllByCondition(ctx, apiInfos, func(tx *gorm.DB) *gorm.DB {
		return tx
	})
	return apiInfos, err
}

func (as *APIServiceImpl) UpdateAPIInfo(ctx context.Context, apiInfo model.APIInfo) error {
	return as.baseService.UpdateById(ctx, &apiInfo)
}

func (as *APIServiceImpl) UpdateAPIInfoByName(ctx context.Context, apiInfo model.APIInfo, name string) error {
	return as.baseService.UpdateByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *APIServiceImpl) DeleteAPIInfo(ctx context.Context, name string) error {
	var apiInfo model.APIInfo
	return as.baseService.DeleteByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name = ?", name)
	})
}

func (as *APIServiceImpl) GetAPIInfoById(ctx context.Context, id uint) (model.APIInfo, error) {
	var apiInfo model.APIInfo
	err := as.baseService.GetByCondition(ctx, &model.APIInfo{}, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
	return apiInfo, err
}

func (as *APIServiceImpl) AddAPIs(ctx context.Context, apiInfos []*model.APIInfo) error {
	return as.baseService.CreateBatch(ctx, apiInfos)
}

func (as *APIServiceImpl) UpdateAPIInfoById(ctx context.Context, apiInfo model.APIInfo, id uint) error {
	return as.baseService.UpdateByCondition(ctx, &apiInfo, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}

func (as *APIServiceImpl) DeleteAPIInfoById(ctx context.Context, id uint) error {
	return as.baseService.DeleteByCondition(ctx, &model.APIInfo{}, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	})
}
