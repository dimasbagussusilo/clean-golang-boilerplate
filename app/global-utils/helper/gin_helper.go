package helper

import (
	"appsku-golang/app/global-utils/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, body model.Response) {
	if body.StatusCode >= 100 && body.StatusCode < 200 {
		body.StatusMessage = "information"
	} else if body.StatusCode >= 200 && body.StatusCode < 300 {
		body.StatusMessage = "success"
	} else if body.StatusCode >= 300 && body.StatusCode < 400 {
		body.StatusMessage = "redirect"
	} else if body.StatusCode >= 400 && body.StatusCode < 500 {
		body.StatusMessage = "client error"
		if body.StatusCode == http.StatusNotFound {
			body.StatusMessage = "success"
			body.StatusCode = http.StatusOK
			body.Data = nil
			body.Error = &model.ErrorLog{
				Message:       "OK",
				SystemMessage: "Data tidak ditemukan",
			}
			body.Total = 0
		}
	} else if body.StatusCode >= 500 && body.StatusCode < 600 {
		body.StatusMessage = "internal server error"
	} else {
		body.StatusMessage = "failed"
	}

	c.JSON(body.StatusCode, body)
}
