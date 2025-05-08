package kafka

import (
	kafka "order-service/controllers/kafka/payment"
	"order-service/services"
)

type Registry struct {
	services services.IServiceRegistry
}

type IKafkaRegistry interface {
	GetPayment() kafka.IPaymentKafka
}

func NewKafkaRegistry(services services.IServiceRegistry) IKafkaRegistry {
	return &Registry{services: services}
}

func (r *Registry) GetPayment() kafka.IPaymentKafka {
	return kafka.NewPaymentKafka(r.services)
}
