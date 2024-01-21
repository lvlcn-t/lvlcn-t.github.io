package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
)

type HandlerConstructor func(data *config.Data) Handler

type Handler interface {
	View(c *gin.Context)
}
