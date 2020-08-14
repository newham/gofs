package api

import "testing"

func TestGetFolder(t *testing.T) {
	SetRoot("../")
	folder := GetFolder("/", nil)
	folder.Print()
}

func TestMkdir(t *testing.T) {
	Mkdir(ROOT_PATH + "new")
}
