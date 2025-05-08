package controllers

import (
	controllers "order-service/controllers/http/order"
	"order-service/services"
)

type Registry struct {
	services services.IServiceRegistry
}

type IControllersRegistry interface {
	GetOrder() controllers.IOrderControllers
}

func NewControllerRegistry(services services.IServiceRegistry) IControllersRegistry {
	return &Registry{services: services}
}

func (r *Registry) GetOrder() controllers.IOrderControllers {
	return controllers.NewOrderControllers(r.services)
}
