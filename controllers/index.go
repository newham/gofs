package controllers

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/newham/gofs/api"
	"github.com/newham/hamgo"
)

func FolderController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) { //get请求页面
		path := getPath(ctx, "/folder/")
		folders := api.GetFolder(path, nil)
		hamgo.Log.Info("%d %-10s %s", 200, "folder", path)
		ctx.JSONHTML(folders, "public/index.html")
		return
	}).OnPOST(func(ctx hamgo.Context) { //post请求json
		m, err := ctx.BindMap()
		if err != nil {
			ctx.JSONMsg(500, "error", err.Error())
		}
		folders := api.GetFolder(m["dir"].(string), nil)
		ctx.JSONFrom(200, folders)
		return
	}).OnPUT(func(ctx hamgo.Context) { //put请求json新建
		m, err := ctx.BindMap()
		if err != nil {
			ctx.JSONMsg(500, "msg", err.Error())
			return
		}
		path := m["dir"].(string) //这是由用户输入的，不用转base64
		filePath := api.ROOT_PATH + path
		// println("path", filePath)
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			ctx.JSONMsg(500, "msg", err.Error())
		} else {
			ctx.JSONOk()
		}
	})
}

func getPath(ctx hamgo.Context, prefix string) string {
	path := strings.TrimPrefix(ctx.R().URL.String(), prefix)
	// if isUnicode {
	// 	path, _ = url.QueryUnescape(strings.TrimPrefix(ctx.R().URL.String(), prefix))
	// }
	path = api.Base64ToURL(path)
	if path == "" {
		path = "/"
	}
	return path
}

func FileController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) {
		path := getPath(ctx, "/file/")
		hamgo.Log.Info("%d %-10s %s", 200, "file", path)
		ctx.File(api.ROOT_PATH + path)
		return
	}).OnDELETE(func(ctx hamgo.Context) {
		deleteMap := map[string]string{}
		err := ctx.BindJSON(&deleteMap)
		if err != nil {
			ctx.JSONMsg(400, "error", err.Error())
			return
		}
		err = api.DeleteFiles(deleteMap)
		if err != nil {
			hamgo.Log.Error("%d %-10s %s", 500, "delete ", "failed:"+err.Error())
			ctx.JSONMsg(500, "error", err.Error())
			return
		}
		ctx.JSONOk()
		return
	}).OnPUT(func(ctx hamgo.Context) {
		path := getPath(ctx, "/file/")
		hamgo.Log.Info("%d %-10s %s", 200, "put file", path)
		_, err := os.Create(api.ROOT_PATH + path)
		if err != nil {
			ctx.JSONMsg(500, "error", err.Error())
			return
		}
		ctx.JSONOk()
	})
}

func IndexController(ctx hamgo.Context) {
	ctx.Redirect("/folder/")
}

func VideoController(ctx hamgo.Context) {
	video := getPath(ctx, "/video/")
	hamgo.Log.Info("%d %-10s %s", 200, "video", video)
	ctx.PutData("video", api.URLToBase64(video))
	ctx.PutData("type", api.GetType(video))
	ctx.PutData("playList", hamgo.JSONToString(api.GetFolder(path.Dir(video), []string{"flv", "video"})))
	ctx.HTML("public/player.html")
}

func UploadController(ctx hamgo.Context) {
	ctx.R().ParseMultipartForm(32 << 20)
	path := ctx.FormValue("path")
	// println("path:", path)
	//1.get upload file
	file, handle, err := ctx.FormFile("file")
	if err != nil {
		// getHtml("").Execute(w, CommonResponse{getMsg("Upload Failed!"), getFolder(path)})
		hamgo.Log.Error("%d %-10s %s", 500, "upload ", "failed:"+err.Error())
		ctx.JSONMsg(http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer file.Close()
	//2.create local file
	f, err := os.OpenFile(api.ROOT_PATH+path+handle.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		// getHtml("").Execute(w, CommonResponse{getMsg("Upload Failed!"), getFolder(path)})
		hamgo.Log.Error("%d %-10s %s", 500, "upload ", handle.Filename+" failed:"+err.Error())
		ctx.JSONMsg(http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer f.Close()
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	hamgo.Log.Info("%d %-10s %s", 200, "upload", path+handle.Filename)
	// getHtml("").Execute(w, CommonResponse{"success", getFolder(path)})
	ctx.JSONMsg(http.StatusOK, "success", true)
}
