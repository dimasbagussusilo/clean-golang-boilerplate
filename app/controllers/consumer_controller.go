package controllers

import (
	"appsku-golang/app/models"
	"appsku-golang/app/usecases"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"appsku-golang/app/global-utils/helper"

	"github.com/sirupsen/logrus"
)

type ConsumerController struct {
	ConsumerUseCase usecases.IConsumerUseCase
}

func NewConsumerController(usecase usecases.IConsumerUseCase) *ConsumerController {
	return &ConsumerController{
		ConsumerUseCase: usecase,
	}
}

func (c *ConsumerController) WorkerUpsertTargetSalesman(ctx context.Context, m kafka.Message) bool {
	var message models.Example

	logrus.WithField(helper.GetRequestIDContext(ctx)).Infof("Processing message from kafka")

	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		logrus.WithField(helper.GetRequestIDContext(ctx)).Error(err)
		return false
	}

	err = c.ConsumerUseCase.ProcessDataUpsert(ctx, message)
	if err != nil {
		logrus.WithField(helper.GetRequestIDContext(ctx)).Error(err)
		return false
	}

	logrus.WithField(helper.GetRequestIDContext(ctx)).Infof("Successful process message from kafka")

	return true
}
