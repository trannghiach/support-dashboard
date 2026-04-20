package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error struct {
		Code string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func JSONError(c *gin.Context, statusCode int, code string, message string) {
	var resp ErrorResponse
	resp.Error.Code = code
	resp.Error.Message = message
	
	c.JSON(statusCode, resp)
}