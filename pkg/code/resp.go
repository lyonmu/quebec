package code

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code,omitempty" example:"50000"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty" example:"success"`
}

func (r Response) Error() string {
	return r.Message
}

func (r Response) Success(data any, c *gin.Context) {
	r.Data = data
	c.JSON(http.StatusOK, r)
}

func (r Response) Failed(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func (r Response) Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, r)
}

func (r Response) Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, r)
}

func (r Response) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, r)
}