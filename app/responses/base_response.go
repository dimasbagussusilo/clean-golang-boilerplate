package responses

import (
	"net/http"

	"appsku-golang/app/global-utils/model"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, body model.Response) {
	body.StatusMessage = "success"
	if body.StatusCode != http.StatusOK && body.StatusCode != http.StatusCreated {
		body.StatusMessage = "failed"
	}

	c.JSON(body.StatusCode, body)
}
