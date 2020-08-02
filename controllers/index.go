package controllers

import (
	"archive/zip"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/newham/gofs/api"
)

const (
	MSG_NONE                   = 0
	MSG_SUCCESS                = 1
	MSG_ERROR                  = 2
	EXTENTION_NULL             = 0
	EXTENTION_EDIT             = 1
	HEADER_CONTENT_DISPOSITION = "Content-Disposition"
	USER_FILE                  = "USER"
)

var ROOT_PATH = api.AppConfig.String("http_file_path")

var SESSION_FILE api.Config

var EDITABLE_TYPE = []string{"txt", "md", "markdown", "h", "c", "cpp", "c++", "go", "xml", "json", "java", "conf", "ini", "css", "js", "sh", "py", "log"}

const HOME = "home/"

type Msg struct {
	Text string
	Type int
}

type CommonResponse struct {
	Msg      Msg
	Folder   Folder
	Username string
}

type SearchResponse struct {
	Msg     Msg
	Results []File
}
type Readme struct {
	Text     string
	FileName string
}
type Folder struct {
	Up      string
	Path    string
	Folders []string
	Files   []File
	Readme  Readme
}
type File struct {
	Name              string
	Size              string
	Path              string
	ModTime           string
	Type              string
	Editable          bool
	DownloadFrequency int
}
type EditResponse struct {
	FileName string
	FilePath string
	FileText string
}

func init() {
	initRoot()
}

func initRoot() {
	if !checkFileIsExist(ROOT_PATH) {
		err := os.MkdirAll(ROOT_PATH, 0777)
		if err != nil {
			log(500, "init", "cant not create "+ROOT_PATH)
			panic(err)
		}
	}
}

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	//read config,need_login if false,do not check session
	if !api.AppConfig.DefaultBool("need_login", false) {
		return true
	}
	//auth by token
	if checkToken(r) {
		return true
	}
	//check session
	if api.HasSession(r) {
		return true
	} else {
		w.WriteHeader(401)
		err := getHtml("login", "").Execute(w, map[string]string{"Msg": "ERROR:Login first!", "Type": "text-danger"})
		if err != nil {
			panic(err)
		}

		log(401, "login", "login.html")
		return false
	}
}

func LogoutController(w http.ResponseWriter, r *http.Request) {
	api.DeleteSession(w, r)
	err := getHtml("login", "").Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

//func ViewController(w http.ResponseWriter, r *http.Request) {
//	t := r.FormValue("t")
//	println(t)
//	session := api.GetSession(r)
//	session.SetView(t)
//	session.Update()
//	log(200, "view", t)
//	redirectPath(w, r, getHome(session.GetUsername()))
//}

func RegisterController(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	token := r.FormValue("token")
	var msg interface{}
	var code int
	if username == "" || password == "" || confirmPassword != password || token != api.AppConfig.String("token") {
		code = 400
		msg = map[string]string{"Msg": "ERROR:username or password wrong!", "Type": "text-danger"}
	} else {
		err := api.AppendString(USER_FILE, fmt.Sprintf("%s,%s\n", username, password))
		if err != nil {
			code = 500
			msg = map[string]string{"Msg": "ERROR:save user failed!", "Type": "text-danger"}
		} else {
			code = 200
			msg = map[string]string{"Msg": "SUCCESS:register", "Type": "text-success"}
		}
	}
	toLogin(w, code, msg)
}

func toLogin(w http.ResponseWriter, code int, msg interface{}) {
	err := getHtml("login", "").Execute(w, msg)
	if err != nil {
		panic(err)
	}
	log(code, "login", "login.html")
}

const (
	VIEW_GRID = "grid"
	VIEW_TBL  = "table"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	//check pwd
	if checkUserPwd(username, password) {
		api.SetSession(api.NewSession(username, VIEW_TBL), w)
		mkhome(username)
		//HttpController(w, r, username)
		redirectPath(w, r, getHome(username))
		return
	}
	//set msg
	var msg interface{}
	var code int
	if r.Method == http.MethodPost {
		msg = map[string]string{"Msg": "ERROR:wrong username or password!", "Type": "text-danger"}
		code = 400

	} else {
		msg = nil
		code = 200
	}
	//return html
	toLogin(w, code, msg)
}

func mkhome(username string) {
	path := ROOT_PATH + getHome(username)
	if !api.IsFileExist(path) {
		os.MkdirAll(path, 0777)
	}
}

func checkUserPwd(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	str, err := api.ReadString(USER_FILE)
	if err != nil {
		return false
	}
	if strings.Contains(str, fmt.Sprintf("%s,%s", username, password)) {
		return true
	}
	return false
}

func IndexController(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view/index.html", "view/body.html", "view/head.html")
	if err != nil {
		panic(err)
	}
	//1.
	err = t.Execute(w, map[string]string{"Msg": "Hello World"})
	//2.
	// err = t.Execute(w, &CommonResponse{"Hello World"})
	if err != nil {
		panic(err)
	}
}

func AboutController(w http.ResponseWriter, r *http.Request) {
	f, _ := os.OpenFile("LICENSE", os.O_RDONLY, 0777)
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	err := getHtml("about", "").Execute(w, map[string]string{"License": string(b)})
	//2.
	// err = t.Execute(w, &CommonResponse{"Hello World"})
	if err != nil {
		panic(err)
	}

	log(200, "about", "")
}

func HttpController(w http.ResponseWriter, r *http.Request, username string) {
	if username == "" {
		username = api.GetUsername(r)
	}
	//1.
	err := getHtml("", api.GetSession(r).GetView()).Execute(w, CommonResponse{getMsg(""), getFolder(getHome(username)), username})
	//2.
	// err = t.Execute(w, &CommonResponse{"Hello World"})
	if err != nil {
		panic(err)
	}

	log(200, "folder", "/")
}

func DownloadShareController(w http.ResponseWriter, r *http.Request) {
	shareKey := r.FormValue("shareKey")
	if shareKey != "" && shareMap[shareKey] != "" {
		serveFile(w, r, shareMap[shareKey])
		log(200, "shareKey", shareKey+":"+shareMap[shareKey])
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("404 Not Found"))
	log(404, "shareKey", shareKey+":"+shareMap[shareKey])
}

func DownloadController(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("name")
	fileType := r.FormValue("type")

	if !checkPermission(w, r, fileName) {
		return
	}
	if fileType == "file" {
		serveFile(w, r, fileName)
		log(200, "file", fileName)
	} else if fileType == "folder" {

	}
}

func serveFile(w http.ResponseWriter, r *http.Request, fileName string) {
	addDownloadFrequency(fileName)
	w.Header().Set(HEADER_CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", getFileName(fileName)))
	http.ServeFile(w, r, ROOT_PATH+fileName)
}

func UploadController(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	path := r.FormValue("filePath")
	// println("path:", path)
	//1.get upload file
	file, handle, err := r.FormFile("file")
	if err != nil {
		// getHtml("").Execute(w, CommonResponse{getMsg("Upload Failed!"), getFolder(path)})
		log(500, "upload ", "failed:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("{\"status\":\"failed\"}"))
		w.Write([]byte("{\"error\":\"failed\"}"))
		return
	}
	defer file.Close()
	//2.create local file
	f, err := os.OpenFile(ROOT_PATH+path+handle.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		// getHtml("").Execute(w, CommonResponse{getMsg("Upload Failed!"), getFolder(path)})
		log(500, "upload ", handle.Filename+" failed:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("{\"status\":\"failed\"}"))
		w.Write([]byte("{\"error\":\"failed\"}"))
		return
	}
	defer f.Close()
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	log(200, "upload", path+handle.Filename)
	// getHtml("").Execute(w, CommonResponse{"success", getFolder(path)})
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("{\"status\":\"success\"}"))
	w.Write([]byte("{\"success\":true}"))
}

func DelController(w http.ResponseWriter, r *http.Request) {
	delType := r.FormValue("type")
	if delType == "array" {
		array := strings.Split(r.FormValue("array")[1:], "|")
		fileName := array[0]
		currentPath := getCurrentDirectory(fileName)
		if strings.HasSuffix(fileName, "/") {
			currentPath = getParentDirectory(fileName)
		}
		if !checkPermission(w, r, currentPath) {
			return
		}
		for _, v := range array {
			if err := os.RemoveAll(ROOT_PATH + v); err != nil {
				panic(err)
			}
		}
		//getHtml("").Execute(w, CommonResponse{getMsg("Delete [" + "array" + "] Success"), getFolder(currentPath)})
		redirectPath(w, r, currentPath)
		log(200, "rm", strings.Join(array, ","))
	} else {
		fileName := r.FormValue("name")
		var currentPath string
		if err := os.RemoveAll(ROOT_PATH + fileName); err != nil {
			panic(err)
		}
		if strings.HasSuffix(fileName, "/") {
			currentPath = getParentDirectory(fileName)
		} else {
			currentPath = getCurrentDirectory(fileName)
		}
		if !checkPermission(w, r, currentPath) {
			return
		}
		//getHtml("").Execute(w, CommonResponse{getMsg("Delete [" + fileName + "] Success"), getFolder(currentPath)})
		redirectPath(w, r, currentPath)
		log(200, "rm", fileName)
	}
}

func redirectPath(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, fmt.Sprintf("/folder?name=%s", path), 301)
}

func FileController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fileName := r.FormValue("name")
		path := r.FormValue("path")
		filePath := ROOT_PATH + path + fileName

		if !checkPermission(w, r, path+fileName) {
			return
		}

		cmd := "touch"

		code, msg := fileBash(cmd, filePath)

		redirectPath(w, r, path)

		log(code, cmd, msg)
	}
}

func FolderController(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		folderName := r.FormValue("name")
		format := r.FormValue("format")
		t := r.FormValue("t")
		session := api.GetSession(r)
		if t == "" {
			t = session.GetView()
		} else {
			session.SetView(t)
			session.Update()
		}

		if !checkPermission(w, r, folderName) {
			return
		}

		if format == "json" {
			b, _ := json.Marshal(CommonResponse{getMsg(""), getFolder(folderName), session.GetUsername()})
			w.WriteHeader(200)
			w.Write(b)
		} else {
			getHtml("", t).Execute(w, CommonResponse{getMsg(""), getFolder(folderName), session.GetUsername()})

		}

		log(200, "folder", folderName)
	} else if r.Method == http.MethodPost {
		folderName := r.FormValue("name")
		path := r.FormValue("path")
		filePath := ROOT_PATH + path + folderName

		if !checkPermission(w, r, path+folderName) {
			return
		}

		cmd := "mkdir"

		code, msg := fileBash(cmd, filePath)

		redirectPath(w, r, path)

		log(code, cmd, msg)
	}

}

func fileBash(cmd string, filePath string) (int, string) {
	var msg string
	var code int
	var err error
	switch cmd {
	case "mkdir":
		err = os.MkdirAll(filePath, 0777)
	case "touch":
		if checkFileIsExist(filePath) {
			err = errors.New("")
		} else {
			_, err = os.Create(filePath)
		}

	case "rm":
		if !checkFileIsExist(filePath) {
			err = errors.New("")
		} else {
			err = os.RemoveAll(filePath)
		}
	default:
		err = errors.New("")
	}
	if err != nil {
		msg = " error:" + err.Error()
		code = 400
	} else {
		msg = ""
		code = 200
	}
	return code, filePath + msg
}

func BashController(w http.ResponseWriter, r *http.Request) {
	var msg string
	var code int
	var err error
	path := r.FormValue("path")
	shell := strings.TrimSpace(r.FormValue("shell"))
	index := strings.Index(shell, " ")
	if index < 0 {
		err = errors.New("")
	} else {

		key := strings.TrimSpace(shell[:index])
		value := strings.TrimSpace(shell[index+1:])

		filePath := ROOT_PATH + path + value
		switch key {
		case "mkdir":
			err = os.MkdirAll(filePath, 0777)
		case "touch":
			if checkFileIsExist(filePath) {
				err = errors.New("")
			} else {
				_, err = os.Create(filePath)
			}

		case "rm":
			if !checkFileIsExist(filePath) {
				err = errors.New("")
			} else {
				err = os.RemoveAll(filePath)
			}
		default:
			err = errors.New("")
		}
	}

	if err != nil {
		msg = " Failed"
		code = 400
	} else {
		msg = " Success"
		code = 200
	}

	getHtml("", "").Execute(w, CommonResponse{getMsg(shell + msg), getFolder(path), api.GetUsername(r)})

	log(code, "bash", shell+msg)
}

func SearchController(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("key")
	path := getHome(api.GetUsername(r))
	if api.GetUsername(r) == "admin" {
		path = ""
	}
	result := search(path, fileName)
	if len(result) > 0 {
		getHtml("search", "").Execute(w, SearchResponse{getMsg("Search [" + fileName + "] Success"), result})
		log(200, "search", fileName)
	} else {
		getHtml("search", "").Execute(w, SearchResponse{getMsg("Search [" + fileName + "] Failed"), result})
		log(403, "search", fileName)
	}

}

func EditController(w http.ResponseWriter, r *http.Request) {
	editType := r.FormValue("type")
	fileName := r.FormValue("name")
	if editType == "open" {
		getHtml("edit", "").Execute(w, EditResponse{fileName, fileName, readFile(ROOT_PATH + fileName)})
	} else if editType == "save" {
		fileText := r.FormValue("text")
		var msg string
		if writeFile(fileName, fileText) {
			msg = "Edit [" + fileName + "] Success"
		} else {
			msg = "Edit [" + fileName + "] Failed"
		}
		getHtml("", "").Execute(w, CommonResponse{getMsg(msg), getFolder(getCurrentDirectory(fileName)), api.GetUsername(r)})
	}
}

var shareMap = map[string]string{}
var shareIndexMap = map[string]string{}

func ShareController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileName := r.FormValue("name")

		if !checkPermission(w, r, fileName) {
			return
		}
		shareKey := shareIndexMap[fileName]
		if shareKey == "" {
			shareKey = api.GetUUID()
			shareMap[shareKey] = fileName
			shareIndexMap[fileName] = shareKey
		}

		w.Header().Set("content-type", "application/json")
		b, _ := json.Marshal(map[string]string{
			"file":     fileName,
			"shareKey": shareKey,
		})
		w.Write(b)

		log(200, "share", fileName)
	}
}

func checkPermission(w http.ResponseWriter, r *http.Request, path string) bool {
	username := api.GetUsername(r)
	if username == "admin" || !api.AppConfig.DefaultBool("need_login", false) {
		return true
	}
	res := strings.HasPrefix(path, getHome(username))
	if !res {
		w.WriteHeader(http.StatusUnauthorized)
		getHtml("", api.GetSession(r).GetView()).Execute(w, CommonResponse{getMsg("open failed , permission denied"), getFolder(getHome(api.GetUsername(r))), username})
		log(http.StatusUnauthorized, "folder", path)
	}
	return res
}

func getHome(username string) string {
	home := "/" + HOME + username
	if !strings.HasSuffix(home, "/") {
		home += "/"
	}
	return home
}

func getFileName(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

func getFolder(path string) Folder {
	dir, err := ioutil.ReadDir(ROOT_PATH + path)
	if err != nil {
		log(500, "getFolder", "open "+path+" failed")
		initRoot()
		return Folder{Path: "/"}
	}
	folders := make([]string, 0, 10)
	files := make([]File, 0, 10)
	readme := Readme{"", ""}
	for _, fi := range dir {
		if fi.IsDir() {
			folders = append(folders, fi.Name()+"/")
		} else {
			if strings.ToLower(fi.Name()) == "readme.md" {
				fullFileName := ROOT_PATH + path + fi.Name()
				f, _ := os.OpenFile(fullFileName, os.O_RDONLY, 0777)
				defer f.Close()
				b, _ := ioutil.ReadAll(f)
				readme = Readme{string(b), fi.Name()}
			}
			fileType := getType(fi.Name())
			files = append(files, File{fi.Name(), formatSize(fi.Size()), path, fi.ModTime().String()[:16], fileType, isEditable(fileType), getDownloadFrequency(path + fi.Name())})
		}

	}
	return Folder{getParentDirectory(path), path, folders, files, readme}
}
func isEditable(fileType string) bool {
	for _, v := range EDITABLE_TYPE {
		if v == fileType {
			return true
		}
	}
	return false
}
func getHtml(name string, t string) *template.Template {
	switch name {
	case "search":
		t, err := template.ParseFiles("view/result.html", "view/head.html", "view/nav2.html", "view/search.html", "view/msg.html")
		if err != nil {
			panic(err)
		}
		return t
	case "edit":
		t, err := template.ParseFiles("view/edit.html", "view/head.html")
		if err != nil {
			panic(err)
		}
		return t
	case "about":
		t, err := template.ParseFiles("view/about.html", "view/head.html", "view/nav2.html")
		if err != nil {
			panic(err)
		}
		return t
	case "login":
		t, err := template.ParseFiles("view/login.html", "view/head.html")
		if err != nil {
			panic(err)
		}
		return t
	default:
		v := "view/folder.html"
		if t == VIEW_GRID {
			v = "view/folder-gird.html"
		}
		t, err := template.ParseFiles("view/index.html", "view/head.html", "view/nav.html", v, "view/msg.html", "view/markdown.html", "view/bar.html")
		if err != nil {
			panic(err)
		}
		return t
	}

}

func getParentDirectory(dirctory string) string {
	p := path.Dir(path.Dir(dirctory))
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}
func getCurrentDirectory(file string) string {
	p := path.Dir(file)
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}
func getFileDirectory(file string) string {
	p := path.Dir(file)
	// println("file=",file,"p=",p)
	p = strings.Replace(p, ROOT_PATH, "", -1)
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}
func readFile(filePath string) string {
	f, _ := os.OpenFile(filePath, os.O_RDONLY, 0777)
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	return string(b)
}
func writeFile(filePath, text string) bool {
	f, err := os.Create(ROOT_PATH + filePath)
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		return false
	} else {
		return true
	}
}
func getMsg(msgText string) Msg {
	msgText = strings.ToLower(msgText)
	var msgType int
	if strings.Contains(msgText, "success") {
		msgType = MSG_SUCCESS
	} else if strings.Contains(msgText, "failed") {
		msgType = MSG_ERROR
	} else {
		msgType = MSG_NONE
	}
	return Msg{msgText, msgType}
}

func checkErr(err error) {
	if err != nil {
		err.Error()
	}
}

func formatSize(size int64) string {
	var K int64 = 1024
	uMap := map[string]int64{"B": 1, "KB": K, "MB": K * K, "GB": K * K * K, "TB": K * K * K * K}
	for k, v := range uMap {
		r := size / v
		if r < K && r > 0 {
			return fmt.Sprintf("%d%s", r, k)
		}
	}
	return "0B"
}

func getType(fileName string) string {
	ext := path.Ext(fileName)
	if len(ext) < 2 {
		return "file"
	}
	return strings.ToLower(ext[1:])
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func log(code int, controller, text string) {
	fmt.Printf("%-4d%-10s%s\n", code, controller, text)
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) ([]byte, error) {
	d, _ := os.Create(dest)
	defer d.Close()
	defer os.Remove(dest)
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return nil, err
		}
	}
	return ioutil.ReadAll(d)
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			defer f.Close()
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func search(path string, key string) []File {
	dirPth := ROOT_PATH + path
	result := make([]File, 0, 10)

	filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		// if err != nil { //忽略错误
		// 	return err
		// }

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if strings.Contains(strings.ToLower(fi.Name()), strings.ToLower(key)) {
			filename := strings.Replace(filename, "\\", "/", -1)
			fileType := getType(filename)
			result = append(result, File{getFileName(filename), formatSize(fi.Size()), getFileDirectory(filename), fi.ModTime().String()[:16], fileType, isEditable(fileType), getDownloadFrequency(filename)})
			return nil
		}

		return nil
	})

	return result
}

var downloadFrequencyPath = "download_frequency.csv"

var mu sync.Mutex

func getDownloadFrequencyCsv() *api.CSV {
	err := api.Mkfile(downloadFrequencyPath)
	if err != nil {
		println(err.Error())
		return nil
	}
	return api.NewCSV(downloadFrequencyPath)
}

func getDownloadFrequency(name string) int {
	name = api.GetMd5String(name)
	csv := getDownloadFrequencyCsv()
	frequencyStrs := csv.Get(name)
	if len(frequencyStrs) < 2 {
		return 0
	}
	frequency, err := strconv.Atoi(frequencyStrs[1])
	if err != nil {
		return 0
	}
	return frequency
}

func addDownloadFrequency(name string) {
	name = api.GetMd5String(name)
	mu.Lock()
	defer mu.Unlock()
	frequency := getDownloadFrequency(name)
	frequency++
	csv := getDownloadFrequencyCsv()
	csv.Put([]string{name, strconv.Itoa(frequency)})
	csv.Save(downloadFrequencyPath)
}
