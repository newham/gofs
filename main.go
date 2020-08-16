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
	server.Static("public")
	server.Handler("/folder/", controllers.FolderController, "POST,GET,PUT")
	server.Handler("/file/", controllers.FileController, "DELETE,GET,PUT")
	server.Handler("/tmp/=tmp", controllers.FileTmpController, "GET")
	server.Post("/file/rename", controllers.FileRenameController)
	server.Post("/file/move", controllers.FileMoveController)
	server.Get("/video/", controllers.VideoController)
	server.Post("/upload", controllers.UploadController)
	server.Post("/download", controllers.DownloadController)
	server.Get("/edit/", controllers.EditController)
	server.Get("/", controllers.IndexController)

	server.Run()
}
