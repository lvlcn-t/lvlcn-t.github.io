package main

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/internal/generator"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/handlers"
	"github.com/lvlcn-t/ChronoTemplify/pkg/renderer"
)

var pages = map[string]string{
	"/":         "public/index.html",
	"/projects": "public/projects/index.html",
}

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

	for route, path := range pages {
		b, err := generator.GeneratePage(r, route)
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			panic(err)
		}
		os.WriteFile(path, b, 0600)
	}

	if err := r.Run(); err != nil {
		os.Exit(1)
	}
}
