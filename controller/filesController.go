package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func FilesController(context *gin.Context)  {
	file, err := context.FormFile("raw")
	if err != nil {
		log.Fatalln(err)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	dir := filepath.Dir(exe)
	if err != nil {
		log.Fatalln(err)
	}
	filename := uuid.New().String()
	uploads := path.Join(dir, "uploads")
	err = os.MkdirAll(uploads, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	fullpath := path.Join("uploads", filename + filepath.Ext(file.Filename))
	fileErr := context.SaveUploadedFile(file, filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatalln(fileErr)
	}
	context.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
}
