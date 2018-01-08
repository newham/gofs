package main

import (
	"github.com/newham/gofs/api"
	_ "github.com/newham/gofs/router"
	"net/http"
)

var SERVER_PORT = api.AppConfig.String("http_port")

func main() {
	println("start server:" + SERVER_PORT)
	//port
	e := http.ListenAndServe(":"+SERVER_PORT, nil)
	if e != nil {
		panic(e)
	}
}
