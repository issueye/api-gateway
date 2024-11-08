package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Condition func(tx *gorm.DB) *gorm.DB

type DataModelInterface interface {
	GetID() uint
}

// BaseService是一个基础服务结构体，包含数据库连接
type BaseService[T DataModelInterface] struct {
	DB *gorm.DB
}

func NewBaseService[T DataModelInterface](md T, db *gorm.DB) BaseService[T] {
	return BaseService[T]{
		DB: db,
	}
}

// GetDB返回数据库连接对象
func (bs *BaseService[T]) GetDB() *gorm.DB {
	return bs.DB
}

// WithTransaction执行给定函数内的数据库事务
func (bs *BaseService[T]) WithTransaction(ctx context.Context, f func(tx *gorm.DB) error) error {
	tx := bs.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := f(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %v", err)
	}
	return tx.Commit().Error
}

// GetById根据给定的ID获取单个记录
func (bs *BaseService[T]) GetById(ctx context.Context, model T) error {
	return bs.DB.WithContext(ctx).Where("id = ?", model.GetID()).First(model).Error
}

// GetByCondition根据给定条件获取单个记录
func (bs *BaseService[T]) GetByCondition(ctx context.Context, condition Condition) (T, error) {
	var model T
	err := condition(bs.DB.WithContext(ctx)).First(model).Error
	return model, err
}

// GetAllByCondition根据给定条件获取所有匹配的记录
func (bs *BaseService[T]) GetAllByCondition(ctx context.Context, condition Condition) ([]T, error) {
	list := make([]T, 0)
	err := condition(bs.DB.WithContext(ctx)).Find(&list).Error
	return list, err
}

// Create添加单个记录
func (bs *BaseService[T]) Create(ctx context.Context, model T) error {
	return bs.DB.WithContext(ctx).Create(model).Error
}

// CreateBatch批量添加记录
func (bs *BaseService[T]) CreateBatch(ctx context.Context, models []T) error {
	return bs.DB.WithContext(ctx).Create(models).Error
}

// UpdateById根据ID更新单个记录
func (bs *BaseService[T]) UpdateById(ctx context.Context, model T) error {
	return bs.DB.WithContext(ctx).Model(model).Where("id = ?", model.GetID()).Updates(model).Error
}

// UpdateByCondition根据条件更新记录
func (bs *BaseService[T]) UpdateByCondition(ctx context.Context, model T, condition Condition) error {
	return condition(bs.DB.WithContext(ctx).Model(model)).Updates(model).Error
}

// DeleteById根据ID删除单个记录
func (bs *BaseService[T]) DeleteById(ctx context.Context, model T) error {
	return bs.DB.WithContext(ctx).Where("id = ?", model.GetID()).Delete(model).Error
}

// DeleteByCondition根据条件删除记录
func (bs *BaseService[T]) DeleteByCondition(ctx context.Context, condition Condition) error {
	var model T
	return condition(bs.DB.WithContext(ctx)).Delete(model).Error
}
