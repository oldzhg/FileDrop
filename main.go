package main

import (
	"chatDemo/server"
	"github.com/zserge/lorca"
)


func main() {
	go func() {
		server.Run()
	}()
	ui, _ := lorca.New("http://localhost:27149/static", "", 1200, 800)

	<-ui.Done()
	ui.Close()
}


