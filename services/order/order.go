package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clientField "order-service/clients/field"
	clientPayment "order-service/clients/payment"
	clientUser "order-service/clients/users"
	"order-service/common/util"
	"order-service/constants"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repositories repositories.IRepositoryRegistry
	clients      clients.IClientRegistry
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.OrderResponse, error)
	GetOrderByUserID(context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(context.Context, *dto.PaymentData) error
}

func NewOrderService(
	repositories repositories.IRepositoryRegistry,
	clients clients.IClientRegistry,
) IOrderService {
	return &OrderService{
		repositories: repositories,
		clients:      clients,
	}
}

func (o *OrderService) GetAllWithPagination(
	ctx context.Context,
	params *dto.OrderRequestParam,
) (*util.PaginationResult, error) {
	orders, total, err := o.repositories.GetOrder().FindAllWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	orderResults := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		user, err := clients.NewClientRegistry().GetUser().GetUserByUUID(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		orderResults = append(orderResults, dto.OrderResponse{
			UUID:      order.UUID,
			Code:      order.Code,
			UserName:  user.Name,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			OrderDate: order.Date,
			CreatedAt: *order.CreatedAt,
			UpdatedAt: *order.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Page:  params.Page,
		Count: total,
		Limit: params.Limit,
		Data:  orderResults,
	}

	responses := util.GeneratePagination(*pagination)

	return &responses, nil
}

func (o *OrderService) GetByUUID(
	ctx context.Context,
	uuid string,
) (*dto.OrderResponse, error) {
	var (
		order *models.Order
		user  *clientUser.UserData
		err   error
	)

	order, err = o.repositories.GetOrder().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user, err = clients.NewClientRegistry().GetUser().GetUserByUUID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:      order.UUID,
		Code:      order.Code,
		UserName:  user.Name,
		Amount:    order.Amount,
		Status:    order.Status.GetStatusString(),
		OrderDate: order.Date,
		CreatedAt: *order.CreatedAt,
		UpdatedAt: *order.UpdatedAt,
	}

	return &response, nil
}

func (o *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	var (
		order []models.Order
		user  = ctx.Value(constants.User).(*clientUser.UserData)
		err   error
	)

	order, err = o.repositories.GetOrder().FindByUserID(ctx, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderLists := make([]dto.OrderByUserIDResponse, 0, len(order))
	for _, item := range order {
		payment, err := o.clients.GetPayment().GetPaymentByUUID(ctx, item.PaymentID)
		if err != nil {
			return nil, err
		}
		orderLists = append(orderLists, dto.OrderByUserIDResponse{
			Code:        item.Code,
			Amount:      util.FormatRupiah(&item.Amount),
			Status:      item.Status.GetStatusString(),
			OrderDate:   item.Date.String(),
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}

	return orderLists, nil
}

func (o *OrderService) Create(ctx context.Context, request *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		order               *models.Order
		txErr, err          error
		user                = ctx.Value(constants.User).(clientUser.UserData)
		field               *clientField.FieldData
		paymentResponse     *clientPayment.PaymentData
		orderFieldSchedules = make([]models.OrderField, 0, len(request.FieldScheduleIDs))
		totalAmount         float64
	)

	for _, fieldID := range request.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldID)
		field, err = o.clients.GetField().GetFieldByUUID(ctx, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedStatus.String() {
			return nil, errOrder.ErrFieldScheduleAlreadyBooked
		}
	}

	err = o.repositories.GetTx().Transaction(func(tx *gorm.DB) error {
		order, txErr = o.repositories.GetOrder().Create(ctx, tx, &models.Order{
			UserID: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})

		if txErr != nil {
			return txErr
		}

		for _, fieldID := range request.FieldScheduleIDs {
			uuidParsed := uuid.MustParse(fieldID)
			orderFieldSchedules = append(orderFieldSchedules, models.OrderField{
				OrderID:         order.ID,
				FieldScheduleID: uuidParsed,
			})
		}

		txErr = o.repositories.GetOrderField().Create(ctx, tx, orderFieldSchedules)
		if err != nil {
			return txErr
		}

		txErr = o.repositories.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(1 * time.Hour)
		description := fmt.Sprintf("Pembayaran sewa %s", field.FieldName)
		paymentResponse, txErr = o.clients.GetPayment().CreatePaymentLink(ctx, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAt,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: dto.CustomerDetail{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetails{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repositories.GetOrder().Update(ctx, tx, &models.Order{
			PaymentID: paymentResponse.UUID,
		}, order.UUID)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		UserName:    user.Name,
		Amount:      totalAmount,
		Status:      order.Status.GetStatusString(),
		PaymentLink: paymentResponse.PaymentLink,
		OrderDate:   order.Date,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}

	return &response, nil
}

func (o *OrderService) mapPaymentStatusToOrder(request *dto.PaymentData) (constants.OrderStatus, *models.Order) {
	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch request.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: request.PaymentID,
			PaidAt:    request.PaidAt,
			Status:    status,
		}
	case constants.ExpirePaymentStatus:
		status = constants.Expired
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			PaidAt:    request.PaidAt,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.Pending
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			PaidAt:    request.PaidAt,
			Status:    status,
		}
	}

	return status, order
}

func (o *OrderService) HandlePayment(ctx context.Context, request *dto.PaymentData) error {
	var (
		txErr, err          error
		order               *models.Order
		orderFieldSchedules []models.OrderField
	)

	status, body := o.mapPaymentStatusToOrder(request)
	err = o.repositories.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repositories.GetOrder().Update(ctx, tx, body, request.OrderID)
		if txErr != nil {
			return txErr
		}

		order, txErr = o.repositories.GetOrder().FindByUUID(ctx, request.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = o.repositories.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: order.ID,
		})

		if request.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repositories.GetOrderField().FindByID(ctx, order.ID)
			if txErr != nil {
				return txErr
			}

			fieldScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				fieldScheduleIDs = append(fieldScheduleIDs, item.FieldScheduleID.String())
			}

			txErr = o.clients.GetField().UpdateStatus(&dto.UpdateFieldScheduleStatusRequest{
				FieldScheduleIDs: fieldScheduleIDs,
			})
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
