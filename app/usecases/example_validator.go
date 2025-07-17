package usecases

type IExampleValidator interface {
}

type ExampleValidator struct {
}

func NewExampleValidator() IExampleValidator {
	return &ExampleValidator{}
}
