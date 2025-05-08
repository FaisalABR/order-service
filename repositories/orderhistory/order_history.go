package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	"order-service/domain/dto"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrderHistoryRepository struct {
	db *gorm.DB
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *dto.OrderHistoryRequest) error
}

func NewOrderHistoryRepository(db *gorm.DB) IOrderHistoryRepository {
	return &OrderHistoryRepository{db: db}
}

func (o *OrderHistoryRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	param *dto.OrderHistoryRequest,
) error {
	request := &models.OrderHistory{
		OrderID: param.OrderID,
		Status:  param.Status,
	}

	err := tx.WithContext(ctx).Create(&request).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSqlQuery)
	}

	return nil
}
