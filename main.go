package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zserge/lorca"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.Writer.WriteString("hello")
		})
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.POST("/api/v1/texts", TextsController)
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
		router.Run(":8080")
	}()
	ui, _ := lorca.New("http://localhost:8080/static", "", 1200, 800)

	<-ui.Done()
	ui.Close()
	//select {
	//
	//}
}

func TextsController(context *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //
		if err != nil {
			log.Fatalln(err)
		}
		dir := filepath.Dir(exe) // 获取当前可执行文件所在目录
		if err != nil {
			log.Fatalln(err)
		}
		filename := uuid.New().String()  // 生成随机文件名
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		context.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}
