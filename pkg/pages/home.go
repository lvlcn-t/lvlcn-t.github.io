package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/components/layout"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/config"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/handlers"
)

var _ handlers.Handler = (*homeHandler)(nil)

type homeHandler struct {
	data *config.Data
}

func NewHandler(data *config.Data) handlers.Handler {
	return &homeHandler{data: data}
}

func (h *homeHandler) View(c *gin.Context) {
	c.HTML(http.StatusOK, "", layout.Layout(Home(h.data)))
}
