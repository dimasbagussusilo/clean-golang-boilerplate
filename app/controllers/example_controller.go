package controllers

import (
	"appsku-golang/app/responses"
	"appsku-golang/app/usecases"

	"appsku-golang/app/global-utils/model"

	"github.com/gin-gonic/gin"
)

type ExampleController struct {
	ExampleUseCase usecases.IExampleUseCase
}

func NewExampleController(usecase usecases.IExampleUseCase) *ExampleController {
	return &ExampleController{
		ExampleUseCase: usecase,
	}
}

func (c *ExampleController) GetExample(r *gin.Context) {
	var result model.Response

	c.ExampleUseCase.Get()
	result.StatusCode = 200

	responses.SendResponse(r, result)
}
