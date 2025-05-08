package controllers

import (
	"net/http"
	errValidation "order-service/common/error"
	"order-service/common/response"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderControllers struct {
	services services.IServiceRegistry
}

type IOrderControllers interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	GetOrderByUserID(*gin.Context)
	Create(*gin.Context)
}

func NewOrderControllers(services services.IServiceRegistry) IOrderControllers {
	return &OrderControllers{services: services}
}

func (o *OrderControllers) GetAllWithPagination(c *gin.Context) {
	var params dto.OrderRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Error:   err,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	orders, err := o.services.GetOrder().GetAllWithPagination(c, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: orders,
		Gin:  c,
	})
}

func (o *OrderControllers) GetByUUID(c *gin.Context) {
	result, err := o.services.GetOrder().GetByUUID(c, c.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (o *OrderControllers) GetOrderByUserID(c *gin.Context) {
	result, err := o.services.GetOrder().GetOrderByUserID(c)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (o *OrderControllers) Create(c *gin.Context) {
	var request dto.OrderRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Error:   err,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := o.services.GetOrder().Create(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: result,
		Gin:  c,
	})
}
