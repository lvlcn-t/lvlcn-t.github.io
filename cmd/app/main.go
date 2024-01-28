package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/internal/generator"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/register"
	"github.com/lvlcn-t/ChronoTemplify/pkg/renderer"
)

const (
	dataFile  = "./data/data.yaml"
	staticDir = "./static"
)

func main() {
	log := slog.Default()
	r := gin.Default()
	r.HTMLRender = renderer.New()
	r.Static("/static", staticDir)

	data, err := config.LoadData(dataFile)
	if err != nil {
		log.Error("Failed to load yaml data", "error", err, "filepath", dataFile)
		panic(err)
	}

	for route, constructor := range register.HandlerMap {
		handler := constructor(&data)
		r.GET(route, handler.View)
	}

	// Generate the static site
	err = generator.GenerateStaticSite(r, "./public", os.DirFS("."))
	if err != nil {
		log.Error("Failed to generate static site", "error", err)
		panic(err)
	}

	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to run http server", "error", err)
		os.Exit(1)
	}
}
