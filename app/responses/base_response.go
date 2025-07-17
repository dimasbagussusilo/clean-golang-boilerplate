package responses

import (
	"net/http"

	"appsku-golang/app/global-utils/model"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

func SendResponse(c *gin.Context, body model.Response) {
	body.StatusMessage = "success"
	if body.StatusCode != http.StatusOK && body.StatusCode != http.StatusCreated {
		body.StatusMessage = "failed"
	}

	c.JSON(body.StatusCode, body)
}

func SendFiberResponse(c *fiber.Ctx, body model.Response) error {
	body.StatusMessage = "success"
	if body.StatusCode != http.StatusOK && body.StatusCode != http.StatusCreated {
		body.StatusMessage = "failed"
	}

	return c.Status(body.StatusCode).JSON(body)
}
