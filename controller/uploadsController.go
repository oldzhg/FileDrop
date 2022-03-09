package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getUploadsDir() (uploads string)  {
exe, err := os.Executable() //
if err != nil {
log.Fatalln(err)
}
dir := filepath.Dir(exe) // 获取当前可执行文件所在目录
uploads = filepath.Join(dir, "uploads")
return
}

func UploadsController(context *gin.Context) {
if path := context.Param("path"); path != "" {
target := filepath.Join(getUploadsDir(), path)
context.Header("Content-Description", "File Transfer")
context.Header("Content-Transfer-Encoding", "binary")
context.Header("Content-Disposition", "attachment; filename="+path)
context.Header("Content-Type", "application/octet-stream")
context.File(target)
} else {
context.Status(http.StatusNotFound)
}
}
