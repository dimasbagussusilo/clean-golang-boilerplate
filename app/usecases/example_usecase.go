package usecases

import "appsku-golang/app/repositories"

type IExampleUseCase interface {
	Get()
}

type ExampleUseCase struct {
	ExampleRepository repositories.IExampleRepository
}

func NewExampleUseCase(ExampleRepository repositories.IExampleRepository) IExampleUseCase {
	return &ExampleUseCase{
		ExampleRepository: ExampleRepository,
	}
}

func (u *ExampleUseCase) Get() {
	u.ExampleRepository.Get()
}
