package routes

import (
	"appsku-golang/app/controllers"
	"appsku-golang/app/handlers"
	"appsku-golang/app/middlewares"
	"context"
	"net/http"

	"appsku-golang/app/global-utils/model"

	"github.com/gofiber/fiber/v2"
)

func NewFiberRoute(ctrlExample *controllers.ExampleController, ctrlStore *controllers.StoreController) *fiber.App {
	ctx := context.Background()

	fiberStoreCtrl := controllers.NewFiberStoreController(ctrlStore)

	f := handlers.MainFiberHandler(ctx)

	f.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(map[string]interface{}{"status": "OK"})
	})

	v1 := f.Group("/v1")
	v1.Use(middlewares.FiberBasicAuthMiddleware())

	storeV1 := v1.Group("/store")
	storeV1.Post("/", fiberStoreCtrl.Insert)
	storeV1.Post("/with-setting", fiberStoreCtrl.InsertWithSetting)
	storeV1.Get("/", fiberStoreCtrl.GetAll)
	storeV1.Get("/:id", fiberStoreCtrl.GetById)
	storeV1.Patch("/:id", fiberStoreCtrl.Update)
	storeV1.Delete("/:id", fiberStoreCtrl.Delete)

	f.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(model.Response{
			StatusCode: http.StatusNotFound,
			Error: &model.ErrorLog{
				Message:       "Not Found",
				SystemMessage: "Not Found",
			},
		})
	})

	return f
}
