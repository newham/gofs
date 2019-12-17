package controllers

import (
	"encoding/base64"
	"net/http"

	"github.com/newham/gofs/api"
)

func PutController(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("file")
	data := r.FormValue("data")
	//decode
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("bad data"))
		return
	}
	path := ROOT_PATH + getHome(getUsernameByToken(r)) + file
	println(path)
	err = api.OverwriteBytes(path, b)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("write file failed:" + err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
