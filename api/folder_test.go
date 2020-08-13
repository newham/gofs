package api

import "testing"

func TestGetFolder(t *testing.T) {
	SetRoot("../")
	folder := GetFolder("/", nil)
	folder.Print()
}
