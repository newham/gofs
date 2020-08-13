package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/newham/hamgo"
)

type Readme struct {
	Text     string
	FileName string
}

type Folder struct {
	Up   string
	Path string
	// Folders [][]string
	Files []File
	// Readme    Readme
	PathArray [][]string
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

func (f File) Print() {
	println(f.Name, f.Size, f.Path, f.ModTime, f.Type, f.Editable, f.DownloadFrequency)
}

func (f Folder) Print() {
	println("Folder:")
	println("UP:", f.Up)
	println("Path:", f.Path)
	println("Files:")
	for _, file := range f.Files {
		file.Print()
	}
	// println("\nFolders:")
	// for _, folder := range f.Folders {
	// 	println(strings.Join(folder, ","))
	// }
}

var DOWNLOAD_FREQUENCY_PATH = "download_frequency.csv"

var ROOT_PATH = "./"

var EDITABLE_TYPE = []string{"txt", "md", "markdown", "h", "c", "cpp", "c++", "go", "xml", "json", "java", "conf", "ini", "css", "js", "sh", "py", "log"}

func SetRoot(path string) {
	if path != "" {
		ROOT_PATH = getPath(path)
	}
	os.MkdirAll(ROOT_PATH, 0777)
}

func GetFolder(path string, typeFilter []string) Folder {
	path = getPath(path)
	dir, err := ioutil.ReadDir(ROOT_PATH + path)
	if err != nil {
		return Folder{Path: "/"}
	}
	// folders := make([][]string, 0, 10)
	files := make([]File, 0, 10)
	// readme := Readme{"", ""} // 这里不要直接读取readme，应该异步请求
	for _, fi := range dir {
		fileType := GetType(fi.Name())
		//过滤不符合typeFilter
		if typeFilter != nil && !strings.Contains(strings.Join(typeFilter, ","), fileType) {
			continue
		}
		//不显示隐藏文件
		if !hamgo.Conf.DefaultBool("show_hidden", false) && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		// if fi.IsDir() {
		// folders = append(folders, []string{fi.Name(), getPath(path + fi.Name())})
		// } else {
		// if strings.ToLower(fi.Name()) == "readme.md" {
		// 	fullFileName := path + fi.Name()
		// 	f, _ := os.OpenFile(fullFileName, os.O_RDONLY, 0444)
		// 	defer f.Close()
		// 	b, _ := ioutil.ReadAll(f)
		// 	readme = Readme{string(b), fi.Name()}
		// }
		filePath := getFile(path + fi.Name())
		if fi.IsDir() {
			fileType = "folder"
			filePath = getPath(filePath)
		}
		files = append(files, File{fi.Name(), formatSize(fi.Size()), filePath, fi.ModTime().String()[:16], fileType, isEditable(fileType), getDownloadFrequency(path + fi.Name())})
		// }

	}
	return Folder{GetParentDirectory(path), path, files, getPathArray(path)}
}

func getPathArray(path string) [][]string {
	// println(path)
	pathArray := [][]string{}
	pathArray = append(pathArray, []string{"全部", "/"}) //根目录
	if path == "" || path == "/" {
		return pathArray
	}
	paths := strings.Split(strings.Trim(path, "/"), "/")
	// println(strings.Join(paths, ","))
	temp := ""
	for _, item := range paths {
		temp += item + "/"
		pathArray = append(pathArray, []string{item, temp})
	}
	return pathArray
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

func isEditable(fileType string) bool {
	for _, v := range EDITABLE_TYPE {
		if v == fileType {
			return true
		}
	}
	return false
}

func GetParentDirectory(dirctory string) string {
	return getPath(path.Dir(path.Dir(dirctory)))
}

func getPath(p string) string {
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	p = strings.ReplaceAll(p, "//", "/")
	if p == "/" || p == "./" {
		p = ""
	}
	return p
}

func getFile(p string) string {
	if strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	p = strings.ReplaceAll(p, "//", "/")
	return p
}

func GetType(fileName string) string {
	ext := path.Ext(fileName)
	extStr := ""
	if len(ext) < 2 {
		extStr = "file"
	} else {
		extStr = strings.ToLower(ext[1:])
	}
	fileType := "nor"
	switch extStr {
	case "txt", "md", "log", "conf", "ini", "plist", "sh", "json", "in", "xml", "css":
		fileType = "txt"
	case "ai":
		fileType = "ai"
	case "psd":
		fileType = "ps"
	case "c", "cpp", "go", "py", "java", "php", "c++":
		fileType = "code"
	case "pdf":
		fileType = "pdf"
	case "png", "jpg", "jpeg", "gif", "bmp", "svg":
		fileType = "pic"
	case "mp4", "mkv", "avi", "rmvb":
		fileType = "video"
	case "dmg":
		fileType = "ipa"
	case "mp3", "wma":
		fileType = "audio"
	case "apk":
		fileType = "apk"
	case "flv":
		fileType = "flv"
	case "doc", "docx":
		fileType = "doc"
	case "xls", "xlsx", "csv":
		fileType = "xls"
	case "ppt", "pptx":
		fileType = "ppt"
	case "tar", "zip", "xz", "gz", "rar", "7z":
		fileType = "zip"
	}
	return fileType
}

func getDownloadFrequency(name string) int {
	name = GetMd5String(name)
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

func getDownloadFrequencyCsv() *CSV {
	err := Mkfile(DOWNLOAD_FREQUENCY_PATH)
	if err != nil {
		println(err.Error())
		return nil
	}
	return NewCSV(DOWNLOAD_FREQUENCY_PATH)
}
