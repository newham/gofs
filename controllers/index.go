package controllers

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
		path := m["dir"].(string)
		if m["base64"].(bool) {
			path = api.Base64ToURL(path)
		}
		folders := api.GetFolder(path, nil)
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
	return api.Base64ToURL(strings.TrimPrefix(ctx.R().URL.String(), prefix))
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
	}).OnPOST(func(ctx hamgo.Context) {
		m, err := ctx.BindMap()
		if err != nil {
			ctx.JSONMsg(500, "error", err.Error())
			return
		}
		path := m["path"].(string)
		txt := m["txt"].(string)
		err = api.OverwriteString(api.ROOT_PATH+path, txt)
		if err != nil {
			ctx.JSONMsg(500, "error", err.Error())
			return
		}
		ctx.JSONOk()
	})
}

func FileTmpController(ctx hamgo.Context) {
	tmp := ctx.PathParam("tmp")
	path := api.TmpDir + tmp
	if !api.IsFileExist(path) {
		ctx.WriteString("404 not found")
		ctx.Text(404)
		return
	}
	ctx.File(path)
}

func FileRenameController(ctx hamgo.Context) {
	m, err := ctx.BindMap()
	if err != nil {
		ctx.JSONMsg(400, "error", err.Error())
		return
	}
	err = api.Rename(m["old"].(string), m["new"].(string), filepath.Dir(api.Base64ToURL(m["path"].(string))))
	if err != nil {
		ctx.JSONMsg(500, "error", err.Error())
		return
	}
	ctx.JSONOk()
}

func FileMoveController(ctx hamgo.Context) {
	req := struct {
		Dir        string
		CheckedMap map[string]string
	}{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSONMsg(400, "error", err.Error())
		return
	}
	err = api.MV(req.CheckedMap, req.Dir)
	if err != nil {
		ctx.JSONMsg(500, "error", err.Error())
		return
	}
	ctx.JSONOk()
}

func IndexController(ctx hamgo.Context) {
	ctx.Redirect("/folder/")
}

func VideoController(ctx hamgo.Context) {
	video := getPath(ctx, "/video/")
	ctx.PutData("playList", hamgo.JSONToString(api.GetFolder(path.Dir(video), func(i int, f api.File) bool {
		if f.Path == api.URLToBase64(video) {
			ctx.PutData("id", i)
		}
		return strings.Contains("video,flv", f.Type)
	})))
	hamgo.Log.Info("%d %-10s %s", 200, "video", video)
	ctx.HTML("public/player.html")
}

func UploadController(ctx hamgo.Context) {
	ctx.R().ParseMultipartForm(32 << 20)
	path := ctx.FormValue("dir")
	//1.get upload file
	file, handle, err := ctx.FormFile("file")
	if err != nil {
		hamgo.Log.Error("%d %-10s %s", 500, "upload ", "failed:"+err.Error())
		ctx.JSONMsg(http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer file.Close()
	//2.create local file
	f, err := os.OpenFile(api.ROOT_PATH+path+handle.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		hamgo.Log.Error("%d %-10s %s", 500, "upload ", handle.Filename+" failed:"+err.Error())
		ctx.JSONMsg(http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer f.Close()
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	hamgo.Log.Info("%d %-10s %s", 200, "upload", path+handle.Filename)
	ctx.JSONMsg(http.StatusOK, "success", true)
}

func DownloadController(ctx hamgo.Context) {
	downloadMap := map[string]string{}
	err := ctx.BindJSON(&downloadMap)
	if err != nil {
		ctx.JSONMsg(400, "error", err.Error())
		return
	}
	fileList := []string{}
	for _, v := range downloadMap {
		fileList = append(fileList, api.ROOT_PATH+api.Base64ToURL(v))
	}
	tmp, err := api.Zip(api.G, fileList)
	if err != nil {
		ctx.JSONMsg(500, "error", err.Error())
		return
	}
	hamgo.Log.Info("%d %-10s %s", 200, "download", tmp)
	ctx.PutData("tmp", tmp)
	ctx.JSONOk()
}

func EditController(ctx hamgo.Context) {
	path := getPath(ctx, "/edit/")
	name := filepath.Base(path)
	txt, err := api.ReadString(api.ROOT_PATH + path)
	if err != nil {
		ctx.Redirect("/")
		return
	}
	ctx.PutData("txt", txt)
	ctx.PutData("name", name)
	ctx.PutData("path", path)
	ctx.HTML("./public/editor.html")
}
