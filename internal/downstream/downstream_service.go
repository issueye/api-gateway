package downstream

import (
	"api-gateway/internal/model"
	"api-gateway/internal/services"
	"context"
)

type DownstreamServiceHandler struct {
	service services.DownstreamServiceImpl
}

func NewDownstreamServiceHandler(service services.DownstreamServiceImpl) *DownstreamServiceHandler {
	return &DownstreamServiceHandler{
		service: service,
	}
}

func (dsh *DownstreamServiceHandler) AddService(ctx context.Context, data *model.Downstream) error {
	return dsh.service.AddService(ctx, data)
}

func (dsh *DownstreamServiceHandler) GetService(ctx context.Context, name string) (model.Downstream, error) {
	return dsh.service.GetService(ctx, name)
}

func (dsh *DownstreamServiceHandler) RemoveService(ctx context.Context, name string) error {
	return dsh.service.RemoveService(ctx, name)
}
