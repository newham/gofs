package controllers

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/newham/gofs/api"
	"github.com/newham/hamgo"
)

func FolderController(ctx hamgo.Context) {
	path := getPath(ctx, "/folder/", true)
	folders := api.GetFolder(path, nil)
	ctx.JSONHTML(folders, "public/index.html")
}

func getPath(ctx hamgo.Context, prefix string, isUnicode bool) string {
	if isUnicode {
		path, _ := url.QueryUnescape(strings.TrimPrefix(ctx.R().URL.String(), prefix))
		return path
	}
	return strings.TrimPrefix(ctx.R().URL.String(), prefix)
}

func FileController(ctx hamgo.Context) {
	ctx.File(api.ROOT_PATH + getPath(ctx, "/file/", true))
}

func IndexController(ctx hamgo.Context) {
	ctx.Redirect("/folder/")
}

func VideoController(ctx hamgo.Context) {
	video := getPath(ctx, "/video/", true)
	ctx.PutData("video", video)
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
		hamgo.Log.Error("%d %s %s", 500, "upload ", "failed:"+err.Error())
		ctx.PutData("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	defer file.Close()
	//2.create local file
	f, err := os.OpenFile(api.ROOT_PATH+path+handle.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		// getHtml("").Execute(w, CommonResponse{getMsg("Upload Failed!"), getFolder(path)})
		hamgo.Log.Error("%d %s %s", 500, "upload ", handle.Filename+" failed:"+err.Error())
		ctx.PutData("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	defer f.Close()
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	hamgo.Log.Info("%d %s %s", 200, "upload", path+handle.Filename)
	// getHtml("").Execute(w, CommonResponse{"success", getFolder(path)})
	ctx.PutData("success", true)
	ctx.JSON(http.StatusOK, nil)
}
