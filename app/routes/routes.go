package routes

import (
	"appsku-golang/app/controllers"
	"appsku-golang/app/handlers"
	"context"
	"net/http"

	"appsku-golang/app/global-utils/model"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	ExampleController *controllers.ExampleController
	StoreController   *controllers.StoreController
}

func NewHttpRoute(ctrlExample *controllers.ExampleController, ctrlStore *controllers.StoreController) *gin.Engine {
	ctx := context.Background()

	routes := Routes{
		ExampleController: ctrlExample,
		StoreController:   ctrlStore,
	}

	g := handlers.MainHttpHandler(ctx)

	g.GET("/health-check", func(context *gin.Context) {
		context.JSON(200, map[string]interface{}{"status": "OK"})
	})

	v1 := g.Group("/v1")
	v1.GET("/example", routes.ExampleController.GetExample)

	storeV1 := v1.Group("/store")
	storeV1.POST("/", routes.StoreController.Insert)
	storeV1.POST("/with-setting", routes.StoreController.InsertWithSetting)
	storeV1.GET("/", routes.StoreController.GetAll)
	storeV1.GET("/:id", routes.StoreController.GetById)
	storeV1.PATCH("/:id", routes.StoreController.Update)
	storeV1.DELETE("/:id", routes.StoreController.Delete)

	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Error: &model.ErrorLog{
				Message:       "Not Found",
				SystemMessage: "Not Found",
			},
		})
	})

	return g
}
