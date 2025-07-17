package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"syscall"
	"time"

	"appsku-golang/app/global-utils/helper"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaPublisher struct {
	Connection          *kafka.Conn
	Brokers             []string
	Writer              *kafka.Writer
	ConsumerGroupReader *kafka.Reader
	Controller          *kafka.Conn
}

type IKafkaPublisher interface {
	GetConnection() *kafka.Conn
	GetController() *kafka.Conn
	CreateTopic(topic string, totalPartition int, totalReplicationFactor int) error
	WriteToTopic(ctx context.Context, topic string, key []byte, message []byte) error
	GetBrokers() []string
	HealthCheck(ctx context.Context) error
}

func NewKafkaPublisher(brokers []string) IKafkaPublisher {
	conn, err := kafka.Dial("tcp", brokers[0])

	if err != nil {
		errStr := "Error failed connect to kafka"
		logrus.Println(errStr)
		logrus.Println(err)
		panic(err)
	}

	controller, err := conn.Controller()

	if err != nil {
		errStr := "Error failed get to controller"
		logrus.Println(errStr)
		logrus.Println(err)
		panic(err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))

	if err != nil {
		errStr := "Error failed connect to controller"
		logrus.Println(errStr)
		logrus.Println(err)
		panic(err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     kafka.CRC32Balancer{},
		BatchTimeout: 5 * time.Millisecond,
		Compression:  kafka.Zstd,
	}

	return &KafkaPublisher{
		Brokers:    brokers,
		Connection: conn,
		Controller: controllerConn,
		Writer:     writer,
	}
}

func (k *KafkaPublisher) GetBrokers() []string {
	return k.Brokers
}

func (k *KafkaPublisher) GetConnection() *kafka.Conn {
	return k.Connection
}

func (k *KafkaPublisher) GetController() *kafka.Conn {
	return k.Controller
}

func (k *KafkaPublisher) CreateTopic(topic string, totalPartition int, totalReplicationFactor int) error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     totalPartition,
			ReplicationFactor: totalReplicationFactor,
		},
	}

	err := k.Controller.CreateTopics(topicConfigs...)

	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			logrus.Println("broken pipe")
		} else {
			errStr := fmt.Sprintf("Error failed create topic %s", topic)
			logrus.Println(errStr)
			logrus.Println(err)
			return err
		}
	}

	return nil
}

func (k *KafkaPublisher) WriteToTopic(ctx context.Context, topic string, key []byte, message []byte) error {

	err := k.Writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Value: message,
			Key:   key,
		},
	)

	if err != nil {
		errStr := fmt.Sprintf("Error failed insert to topic %s", topic)
		logrus.Println(errStr)
		logrus.Println(err)
		return err
	}

	return nil
}

func (k *KafkaPublisher) HealthCheck(ctx context.Context) error {
	testConn, err := kafka.Dial("tcp", net.JoinHostPort(k.Controller.Broker().Host, strconv.Itoa(k.Controller.Broker().Port)))
	if err != nil {
		return fmt.Errorf("failed to connect to controller %s:%d: %w", k.Controller.Broker().Host, k.Controller.Broker().Port, err)
	}
	defer testConn.Close()

	if k.Writer == nil {
		return fmt.Errorf("kafka writer is nil")
	}

	return nil
}

func SetReader(brokers []string, topic string, partition int, offset int64) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3,
		MaxBytes:  10e6,
		MaxWait:   10e6,
	})

	if offset > 0 {
		reader.SetOffset(offset)
	}

	return reader
}

func SetConsumerGroupReader(brokers []string, topic string, groupID string, queueCapacity int) *kafka.Reader {
	kafkaReaderConfig := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}

	if queueCapacity > 0 {
		kafkaReaderConfig.QueueCapacity = queueCapacity
	}

	reader := kafka.NewReader(kafkaReaderConfig)

	return reader
}

type Consumer struct {
	Ctx           context.Context
	Name          string
	Topic         string
	GroupId       string
	Worker        func(ctx context.Context, m kafka.Message) bool
	Partition     int
	Brokers       []string
	Offset        int64
	QueueCapacity int
}

func (c Consumer) ReadFromTopic(ctx context.Context) {
	reader := SetConsumerGroupReader(c.Brokers, c.Topic, c.GroupId, c.QueueCapacity)
	logrus.Infof("Running consumer %s", c.Name)
	for {
		message, err := reader.FetchMessage(c.Ctx)
		if err != nil {
			logrus.WithField(helper.GetRequestIDContext(ctx)).Error(err)
			time.Sleep(30 * time.Second)
			continue
		}

		consumed := c.Worker(c.Ctx, message)
		if !consumed {
			logrus.WithField(helper.GetRequestIDContext(ctx)).Errorf("Failed proccess message %s", string(message.Value))
		}

		// commit message
		err = reader.CommitMessages(c.Ctx, message)
		if err != nil {
			logrus.WithField(helper.GetRequestIDContext(ctx)).Errorf("Error commit message %s", err.Error())
		}
	}

	reader.Close()
}
