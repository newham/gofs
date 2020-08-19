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

func LoginController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) {
		ctx.HTML("public/login.html")
	}).OnPOST(func(ctx hamgo.Context) {
		name := ctx.FormValue("user")
		//check
		if name == "" {
			ctx.PutData("error", "用户名不能为空！")
			//登录失败
			ctx.HTML("public/login.html")
			return
		}

		ctx.FormValue("user")
		user := api.GetUser(api.GetMd5String(name))
		if user != nil && user.PWD == ctx.FormValue("pwd") {
			//设置session
			ctx.GetSession().Set("USER", user)
			ctx.SetCookie("user", name, "/")
			ctx.Redirect("/")
			return
		}
		ctx.PutData("error", "用户不存在/密码错误") //登录失败
		ctx.HTML("public/login.html")
		return
	})
}

func RegisterController(ctx hamgo.Context) {
	// name := ctx.FormValue("user")
	// pwd := ctx.FormValue("pwd")
	// email := ctx.FormValue("email")
	m, err := ctx.BindMap()
	if err != nil {
		ctx.PutData("code", 500)
		respError(ctx, 500, err.Error())
		return
	}

	name := m["user"].(string)
	pwd := m["pwd"].(string)
	email := m["email"].(string)
	// println(name, pwd, email)
	if api.GetUser(api.GetMd5String(name)) != nil {
		ctx.PutData("code", 400)
		respError(ctx, 400, "用户已经存在！")
		return
	}
	api.PutUser(api.NewUser(name, pwd, email, "", ""))
	ctx.JSONOk()
}

func LogoutController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) {
		ctx.DeleteSession()
		ctx.Redirect("/login")
	}).OnPOST(func(ctx hamgo.Context) {
		ctx.JSONOk()
	})
}

func FolderController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) { //get请求页面
		path := getPath(ctx, "/folder/")
		hamgo.Log.Info("%d %-10s %s", 200, "folder", path)
		ctx.HTML("public/index.html")
		return
	}).OnPOST(func(ctx hamgo.Context) { //post请求json
		m, err := ctx.BindMap()
		if err != nil {
			respError(ctx, 400, err.Error())
			return
		}
		path := m["dir"].(string)
		if m["base64"].(bool) {
			path = api.Base64ToURL(path)
		}
		//设置用户目录
		folders := api.GetFolder(getUserRoot(ctx), path, nil)
		ctx.JSONFrom(200, folders)
		return
	}).OnPUT(func(ctx hamgo.Context) { //put请求json新建
		m, err := ctx.BindMap()
		if err != nil {
			respError(ctx, 500, err.Error())
			return
		}
		path := m["dir"].(string) //这是由用户输入的，不用转base64
		filePath := getUserRoot(ctx) + path
		// println("path", filePath)
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			respError(ctx, 500, err.Error())
		} else {
			ctx.JSONOk()
		}
	})
}

func getUser(ctx hamgo.Context) *api.User {
	return ctx.GetSession().Get("USER").(*api.User)
}

func getUserRoot(ctx hamgo.Context) string {
	userName := getUser(ctx).Name
	if userName == "admin" {
		return api.ROOT_PATH
	}
	return api.ROOT_PATH + "home/" + userName + "/"
}

func getUserTmpRoot(ctx hamgo.Context) string {
	return getUserRoot(ctx) + api.TMP_ROOT
}

func getPath(ctx hamgo.Context, prefix string) string {
	return api.Base64ToURL(strings.TrimPrefix(ctx.R().URL.String(), prefix))
}

func FileController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) {
		path := getPath(ctx, "/file/")
		hamgo.Log.Info("%d %-10s %s", 200, "file", path)
		ctx.File(getUserRoot(ctx) + path)
		return
	}).OnDELETE(func(ctx hamgo.Context) {
		deleteMap := map[string]string{}
		err := ctx.BindJSON(&deleteMap)
		if err != nil {
			respError(ctx, 400, err.Error())
			return
		}
		deleteArray := []string{}
		userRoot := getUserRoot(ctx)
		for _, v := range deleteMap {
			deleteArray = append(deleteArray, userRoot+api.Base64ToURL(v))
		}
		err = api.DeleteFiles(deleteArray)
		if err != nil {
			hamgo.Log.Error("%d %-10s %s", 500, "delete ", "failed:"+err.Error())
			respError(ctx, 500, err.Error())
			return
		}
		ctx.JSONOk()
		return
	}).OnPUT(func(ctx hamgo.Context) { //新增，put:{dir,txt}
		m, err := ctx.BindMap()
		if err != nil {
			respError(ctx, 500, err.Error())
			return
		}
		path := getUserRoot(ctx) + m["dir"].(string) //这是由用户输入的，不用转base64
		// println("path", filePath)
		err = api.Mkfile(path)
		if err != nil {
			respError(ctx, 500, err.Error())
		} else {
			ctx.JSONOk()
		}
	}).OnPOST(func(ctx hamgo.Context) { //对已有的修改，post:/=path,{txt}
		path := getPath(ctx, "/file/")
		m, err := ctx.BindMap()
		if err != nil {
			respError(ctx, 500, err.Error())
			return
		}
		txt := m["txt"].(string)
		err = api.OverwriteString(getUserRoot(ctx)+path, txt)
		if err != nil {
			respError(ctx, 500, err.Error())
			return
		}
		ctx.JSONOk()
	})
}

func FileTmpController(ctx hamgo.Context) {
	ctx.OnGET(func(ctx hamgo.Context) {
		tmp := ctx.PathParam("tmp")
		tmpPath := getUserTmpRoot(ctx) + tmp
		if !api.IsFileExist(tmpPath) { //不存在就创建tmp
			// ctx.WriteString("404 not found")
			// ctx.Text(404)
			// return
			api.Mkdir(tmpPath)
		}
		if tmp == "" { //获取tmp列表
			folder := api.GetFolder(getUserRoot(ctx), api.TMP_ROOT, nil)
			ctx.JSONFrom(200, folder)
			return
		}
		ctx.File(tmpPath)
	}).OnDELETE(func(ctx hamgo.Context) {
		tmp := ctx.PathParam("tmp")
		if tmp == "*" {
			api.DeleteFile(getUserTmpRoot(ctx)) //先删除
			api.Mkdir(getUserTmpRoot(ctx))      //再创建一个空的缓存文件夹
			ctx.JSONOk()
			return
		}
	})

}

func FileRenameController(ctx hamgo.Context) {
	m, err := ctx.BindMap()
	if err != nil {
		respError(ctx, 400, err.Error())
		return
	}
	userRoot := getUserRoot(ctx)
	dir := strings.TrimSuffix(api.Base64ToURL(m["path"].(string)), "/")
	path := userRoot + filepath.Dir(dir)
	old := m["old"].(string)
	new := m["new"].(string)
	err = api.Rename(path, old, new)
	if err != nil {
		respError(ctx, 500, err.Error())
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
		respError(ctx, 400, err.Error())
		return
	}
	chechedArray := []string{}
	userRoot := getUserRoot(ctx)
	for _, v := range req.CheckedMap {
		chechedArray = append(chechedArray, userRoot+api.Base64ToURL(v))
	}
	err = api.MV(chechedArray, userRoot+req.Dir)
	if err != nil {
		respError(ctx, 500, err.Error())
		return
	}
	ctx.JSONOk()
}

func IndexController(ctx hamgo.Context) {
	ctx.Redirect("/folder/")
}

func VideoController(ctx hamgo.Context) {
	video := getPath(ctx, "/video/")
	ctx.PutData("playList", hamgo.JSONToString(api.GetFolder(getUserRoot(ctx), path.Dir(video), func(i int, f api.File) bool {
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
		respError(ctx, 500, err.Error())
		return
	}
	defer file.Close()
	//2.create local file
	f, err := os.OpenFile(getUserRoot(ctx)+path+handle.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		hamgo.Log.Error("%d %-10s %s", 500, "upload ", handle.Filename+" failed:"+err.Error())
		respError(ctx, 500, err.Error())
		return
	}
	defer f.Close()
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	hamgo.Log.Info("%d %-10s %s", 200, "upload", path+handle.Filename)
	ctx.JSONMsg(http.StatusOK, "success", true)
}

func DownloadController(ctx hamgo.Context) {
	path := api.Base64ToURL(ctx.PathParam("path")) //需要转码
	filePath := getUserRoot(ctx) + path
	if !api.IsFileExist(filePath) {
		ctx.WriteString("404 not found")
		ctx.Text(404)
		return
	}
	hamgo.Log.Info("%d %-10s %s", 200, "download", path)
	ctx.Attachment(getUserRoot(ctx) + path)
}

func DownloadZipController(ctx hamgo.Context) {
	downloadMap := map[string]string{}
	err := ctx.BindJSON(&downloadMap)
	if err != nil {
		respError(ctx, 400, err.Error())
		return
	}
	fileList := []string{}
	for _, v := range downloadMap {
		fileList = append(fileList, getUserRoot(ctx)+api.Base64ToURL(v))
	}
	tmp, err := api.Zip(getUserRoot(ctx), api.G, fileList)
	if err != nil {
		respError(ctx, 500, err.Error())
		return
	}
	hamgo.Log.Info("%d %-10s %s", 200, "download", tmp)
	ctx.PutData("tmp", tmp)
	ctx.JSONOk()
}

func EditController(ctx hamgo.Context) {
	file := getPath(ctx, "/edit/")
	ctx.PutData("history", hamgo.JSONToString(api.GetFolder(getUserRoot(ctx), path.Dir(file), func(i int, f api.File) bool {
		if f.Path == api.URLToBase64(file) {
			ctx.PutData("id", i)
		}
		if f.Editable {
			return true
		}
		return false
	})))
	ctx.HTML("./public/editor.html")
}

func respError(ctx hamgo.Context, code int, msg string) {
	ctx.JSONMsg(code, "error", msg)
}
