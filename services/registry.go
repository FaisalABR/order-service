package services

import (
	"order-service/clients"
	"order-service/repositories"
	services "order-service/services/order"
)

type Registry struct {
	repositories repositories.IRepositoryRegistry
	clients      clients.IClientRegistry
}

type IServiceRegistry interface {
	GetOrder() services.IOrderService
}

func NewServiceRegistry(
	repositories repositories.IRepositoryRegistry,
	clients clients.IClientRegistry,
) IServiceRegistry {
	return &Registry{repositories: repositories, clients: clients}
}

func (r *Registry) GetOrder() services.IOrderService {
	return services.NewOrderService(r.repositories, r.clients)
}
