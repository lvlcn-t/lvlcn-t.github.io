package generator

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func GeneratePage(r *gin.Engine, route string) ([]byte, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}
	r.ServeHTTP(w, req)
	return w.Body.Bytes(), nil
}
