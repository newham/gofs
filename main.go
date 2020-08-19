package main

import (
	"github.com/newham/gofs/api"
	"github.com/newham/gofs/controllers"
	"github.com/newham/hamgo"
)

func main() {
	//指定配置文件
	server := hamgo.NewByConf("./conf/app.conf")
	//设置下载保存文件的根目录
	api.SetRoot(hamgo.Conf.String("http_file_path"))
	//filter
	server.AddFilter(controllers.SessionFilter).AddAnnoURL("/login").AddAnnoURL("/register")
	//Static
	server.Static("public")
	//controller
	server.Handler("/login", controllers.LoginController, "POST,GET")
	server.Post("/register", controllers.RegisterController)
	server.Handler("/logout", controllers.LogoutController, "POST,GET")
	server.Handler("/folder/", controllers.FolderController, "POST,GET,PUT,POST")
	server.Handler("/file/", controllers.FileController, "DELETE,GET,PUT,POST")
	server.Handler("/tmp/=tmp", controllers.FileTmpController, "GET,DELETE")
	server.Post("/file/rename", controllers.FileRenameController)
	server.Post("/file/move", controllers.FileMoveController)
	server.Get("/video/", controllers.VideoController)
	server.Post("/upload", controllers.UploadController)
	server.Post("/download", controllers.DownloadZipController)
	server.Get("/download/=path", controllers.DownloadController)
	server.Get("/edit/", controllers.EditController)
	server.Get("/", controllers.IndexController)

	server.Run()
}
