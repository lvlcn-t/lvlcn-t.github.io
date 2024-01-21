package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/components/layout"
	"github.com/lvlcn-t/ChronoTemplify/pkg/pages"
)

type project struct{}

func NewProjectsHandler() Handler {
	return &project{}
}

func (p *project) View(c *gin.Context) {
	page := layout.Layout(pages.Projects())
	c.HTML(http.StatusOK, "", page)
}
