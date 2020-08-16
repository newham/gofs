package api

import "testing"

func TestZip(t *testing.T) {
	temp, _ := Zip(1024*1024*1024, []string{"../files", "./zip.go"})
	println(temp)
}
