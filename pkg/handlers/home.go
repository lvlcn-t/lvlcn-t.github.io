package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/components/layout"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/pages"
)

type home struct {
	data *config.Data
}

func NewHomeHandler(data *config.Data) Handler {
	return &home{data: data}
}

func (h *home) View(c *gin.Context) {
	page := layout.Layout(pages.Home(h.data))
	c.HTML(http.StatusOK, "", page)
}
