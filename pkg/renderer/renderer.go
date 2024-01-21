package renderer

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
)

type templRender struct {
	Code int
	Data templ.Component
}

func New() *templRender {
	return &templRender{}
}

func (t templRender) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}
	return nil
}

func (t templRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (t *templRender) Instance(name string, data any) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &templRender{
			Code: http.StatusOK,
			Data: templData,
		}
	}
	return nil
}
