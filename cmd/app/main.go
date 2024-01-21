package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/handlers"
	"github.com/lvlcn-t/ChronoTemplify/pkg/renderer"
)

func main() {
	r := gin.Default()
	r.HTMLRender = renderer.New()

	data, err := config.LoadData("./static/data/data.yaml")
	if err != nil {
		panic(err)
	}

	h := handlers.NewHomeHandler(&data)
	p := handlers.NewProjectsHandler()

	r.Static("/static", "./static")

	_ = r.GET("/", h.View)
	_ = r.GET("/projects", p.View)

	if err := r.Run(); err != nil {
		os.Exit(1)
	}
}
