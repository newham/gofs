package router

import (
	"net/http"
	"strings"

	"github.com/newham/gofs/controllers"
)

func init() {

	//public
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))) //设置css文件目录

	//controllers
	http.HandleFunc("/", BaseController)
}

func BaseController(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	//exception
	if strings.HasPrefix(uri, "/login") {
		controllers.LoginController(w, r)
		return
	}
	if strings.HasPrefix(uri, "/register") {
		controllers.RegisterController(w, r)
		return
	}
	//filter
	if !controllers.CheckSession(w, r) {
		return
	}
	//router
	if strings.HasPrefix(uri, "/folder") {
		controllers.FolderController(w, r)
	} else if strings.HasPrefix(uri, "/file") {
		controllers.FileController(w, r)
	} else if strings.HasPrefix(uri, "/download") {
		controllers.DownloadController(w, r)
	} else if strings.HasPrefix(uri, "/del") {
		controllers.DelController(w, r)
	} else if strings.HasPrefix(uri, "/upload") {
		controllers.UploadController(w, r)
	} else if strings.HasPrefix(uri, "/search") {
		controllers.SearchController(w, r)
	} else if strings.HasPrefix(uri, "/edit") {
		controllers.EditController(w, r)
	} else if strings.HasPrefix(uri, "/about") {
		controllers.AboutController(w, r)
	} else if strings.HasPrefix(uri, "/logout") {
		controllers.LogoutController(w, r)
	} else {
		controllers.HttpController(w, r, "")
	}
	// else if strings.HasPrefix(uri, "/bash") {
	// 	controllers.BashController(w, r)
	// }
}
