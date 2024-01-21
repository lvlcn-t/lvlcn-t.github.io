package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/components/layout"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/handlers"
)

var _ handlers.Handler = (*projectsHandler)(nil)

type projectsHandler struct{}

func NewHandler(_ *config.Data) handlers.Handler {
	return &projectsHandler{}
}

func (p *projectsHandler) View(c *gin.Context) {
	page := layout.Layout(Projects())
	c.HTML(http.StatusOK, "", page)
}
