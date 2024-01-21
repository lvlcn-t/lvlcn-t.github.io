package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	View(c *gin.Context)
}
