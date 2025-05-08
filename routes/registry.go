package routes

import (
	"order-service/clients"
	controllers "order-service/controllers/http"
	routes "order-service/routes/order"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controllers controllers.IControllersRegistry
	clients     clients.IClientRegistry
	group       *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(controllers controllers.IControllersRegistry, clients clients.IClientRegistry, group *gin.RouterGroup) IRouteRegistry {
	return &Registry{
		controllers: controllers,
		clients:     clients,
		group:       group,
	}
}

func (r *Registry) Serve() {
	r.orderRoute().Run()
}

func (r *Registry) orderRoute() routes.IOrderRoutes {
	return routes.NewOrderRoute(r.controllers, r.clients, r.group)
}
