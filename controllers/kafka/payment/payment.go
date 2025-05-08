package kafka

import (
	"context"
	"encoding/json"
	"order-service/common/util"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	services services.IServiceRegistry
}

type IPaymentKafka interface {
	HandlePayment(context.Context, *sarama.ConsumerMessage) error
}

func NewPaymentKafka(services services.IServiceRegistry) IPaymentKafka {
	return &PaymentKafka{services: services}
}

func (p *PaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	defer util.Recover()
	var body dto.PaymentContent
	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return err
	}

	data := body.Body.Data
	err = p.services.GetOrder().HandlePayment(ctx, &data)
	if err != nil {
		logrus.Errorf("failed to handle payment: %v", err)
		return err
	}

	logrus.Infof("success handle payment")
	return nil
}
