package handlers

import (
	"appsku-golang/app/config"
	"appsku-golang/app/constants"
	"appsku-golang/app/controllers"
	"context"

	"appsku-golang/app/global-utils/helper"
	kafkadbo "appsku-golang/app/global-utils/kafka"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type IConsumerHandler interface {
	BindConsumer(ctx context.Context, cfg config.Configuration, queueName string)
}

type ConsumerHandler struct {
	Kafka           kafkadbo.IKafkaPublisher
	Controller      *controllers.ConsumerController
	StoreController any
}

func MainConsumerHandler(ctrl *controllers.ConsumerController, kafkaConn kafkadbo.IKafkaPublisher) IConsumerHandler {
	return &ConsumerHandler{
		Kafka:           kafkaConn,
		Controller:      ctrl,
		StoreController: nil,
	}
}

func (h *ConsumerHandler) BindConsumer(ctx context.Context, cfg config.Configuration, queueName string) {

	consumer := kafkadbo.Consumer{
		Ctx:       helper.SetRequestIDToContext(context.Background(), uuid.NewString()),
		Partition: 1,
		Offset:    0,
		Brokers:   config.Get().MessageBroker.Kafka.Hosts,
	}

	switch queueName {
	case constants.TopicExample:
		consumer.Name = queueName
		consumer.GroupId = consumer.Name
		consumer.Topic = constants.TopicExample
		consumer.Worker = h.Controller.WorkerUpsertTargetSalesman
	case constants.TopicStore:
		consumer.Name = queueName
		consumer.GroupId = consumer.Name
		consumer.Topic = constants.TopicStore
		if h.StoreController != nil {
			consumer.Worker = h.StoreController.(func(ctx context.Context, m kafka.Message) bool)
		} else {
			log.Panicf("StoreController is nil")
			return
		}
	default:
		log.Panicf("Invalid consumer name: %s", consumer.Name)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			// Log the panic and continue running the program
			log.WithField(helper.GetRequestIDContext(ctx)).Errorf("goroutine panicked: %v", r)
		}
	}()

	err := h.Kafka.CreateTopic(consumer.Topic, consumer.Partition, len(h.Kafka.GetBrokers())-1)
	if err != nil {
		log.Fatal(err)
		return
	}

	// create topic for requeue
	err = h.Kafka.CreateTopic("requeue_"+consumer.Topic, consumer.Partition, len(h.Kafka.GetBrokers())-1)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func(consumer kafkadbo.Consumer) {
		consumer.ReadFromTopic(ctx)
	}(consumer)

	<-helper.ExitKafka
}
