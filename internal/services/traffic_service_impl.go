package services

import (
	"api-gateway/internal/global"
	"api-gateway/internal/model"
	"api-gateway/pkg/service"
	"context"
)

type TrafficService struct {
	baseService service.BaseService[*model.TrafficStats]
}

func NewTrafficService() TrafficService {
	bs := service.NewBaseService(&model.TrafficStats{}, global.DB)
	return TrafficService{
		baseService: bs,
	}
}

func (ts *TrafficService) RecordTrafficStats(ctx context.Context, data *model.TrafficStats) error {
	return ts.baseService.Create(ctx, data)
}
