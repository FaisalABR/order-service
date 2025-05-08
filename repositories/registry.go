package repositories

import (
	orderRepository "order-service/repositories/order"
	orderFieldRepository "order-service/repositories/orderfield"
	orderHistoryRepository "order-service/repositories/orderhistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetOrder() orderRepository.IOrderRepository
	GetOrderField() orderFieldRepository.IOrderFieldRepository
	GetOrderHistory() orderHistoryRepository.IOrderHistoryRepository
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetOrder() orderRepository.IOrderRepository {
	return orderRepository.NewOrderRepository(r.db)
}

func (r *Registry) GetOrderField() orderFieldRepository.IOrderFieldRepository {
	return orderFieldRepository.NewOrderFieldRepository(r.db)
}

func (r *Registry) GetOrderHistory() orderHistoryRepository.IOrderHistoryRepository {
	return orderHistoryRepository.NewOrderHistoryRepository(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
