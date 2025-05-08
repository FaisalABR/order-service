package routes

import (
	"order-service/clients"
	"order-service/constants"
	controllers "order-service/controllers/http"
	"order-service/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	controllers controllers.IControllersRegistry
	clients     clients.IClientRegistry
	group       *gin.RouterGroup
}

type IOrderRoutes interface {
	Run()
}

func NewOrderRoute(
	controllers controllers.IControllersRegistry,
	clients clients.IClientRegistry,
	group *gin.RouterGroup,
) IOrderRoutes {
	return &OrderRoutes{
		controllers: controllers,
		clients:     clients,
		group:       group,
	}
}

func (o *OrderRoutes) Run() {
	group := o.group.Group("/orders")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Customer,
		constants.Admin,
	}, o.clients), o.controllers.GetOrder().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Customer,
		constants.Admin,
	}, o.clients), o.controllers.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.clients), o.controllers.GetOrder().GetOrderByUserID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.clients), o.controllers.GetOrder().Create)
}
