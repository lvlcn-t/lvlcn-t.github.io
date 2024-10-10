package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/lvlcn-t.github.io/internal/generator"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/config"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/register"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/renderer"
)

const (
	dataFile  = "./data/data.yaml"
	staticDir = "./static"
)

func main() {
	debug := strings.EqualFold(os.Getenv("DEBUG"), "true")
	log := slog.Default()
	r := gin.New()
	r.Use(gin.Logger())
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.Use(gin.Recovery())
	}

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

	switch debug {
	case false:
		err = generator.GenerateStaticSite(r, "./public", os.DirFS("."))
		if err != nil {
			log.Error("Failed to generate static site", "error", err)
			panic(err)
		}
	case true:
		if err = r.Run(":8080"); err != nil {
			log.Error("Failed to run http server", "error", err)
			os.Exit(1)
		}
	}
}
