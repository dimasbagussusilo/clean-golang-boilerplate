package usecases

import (
	"appsku-golang/app/models"
	"context"
)

type IConsumerUseCase interface {
	ProcessDataUpsert(ctx context.Context, message models.Example) error
}

type ConsumerUseCase struct {
}

func NewConsumerUseCase() IConsumerUseCase {
	return &ConsumerUseCase{}
}

func (u *ConsumerUseCase) ProcessDataUpsert(ctx context.Context, message models.Example) error {

	return nil
}
