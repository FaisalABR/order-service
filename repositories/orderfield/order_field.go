package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrderFieldRepository struct {
	db *gorm.DB
}

type IOrderFieldRepository interface {
	FindByID(context.Context, uint) ([]models.OrderField, error)
	Create(context.Context, *gorm.DB, []models.OrderField) error
}

func NewOrderFieldRepository(db *gorm.DB) IOrderFieldRepository {
	return &OrderFieldRepository{db: db}
}

func (o *OrderFieldRepository) FindByID(ctx context.Context, id uint) ([]models.OrderField, error) {
	var orders []models.OrderField

	err := o.db.WithContext(ctx).Where("id = ?", id).Find(&orders).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSqlQuery)
	}

	return orders, nil
}

func (o *OrderFieldRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	request []models.OrderField,
) error {
	err := tx.WithContext(ctx).Create(&request).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSqlQuery)
	}

	return nil
}
