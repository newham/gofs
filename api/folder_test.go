package api

import "testing"

func TestGetFolder(t *testing.T) {
	SetRoot("../")
	folder := GetFolder(ROOT_PATH, "/", nil)
	folder.Print()
}

func TestMkdir(t *testing.T) {
	Mkdir(ROOT_PATH + "new")
}

func TestDirSize(t *testing.T) {
	println(formatSize(FileSize("../files")))
	println(formatSize(FileSize("./file.go")))
}
