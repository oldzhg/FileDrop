package server

import (
	"FileDrop/controller"
	"FileDrop/server/ws"
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	hub := ws.NewHub()
	go hub.Run()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteString("hello")
	})
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.POST("/api/v1/texts", controller.TextsController)
	router.POST("/api/v1/files", controller.FilesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.StaticFS("/static", http.FS(staticFiles))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatalln(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatalln(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":27149")
}
