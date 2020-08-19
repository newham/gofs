/******************************************
*FileName: file.go
*Author: Liu han
*Date: 2017-11-24
*Description: File Read & Write Tool
*******************************************/
package api

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func ReadBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func ReadString(filename string) (string, error) {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(str), err
}

func AppendBytes(filename string, content []byte) error {
	return writeString(filename, string(content), 1)
}

func AppendString(filename string, content string) error {
	return writeString(filename, content, 1)
}

func OverwriteString(filename string, content string) error {
	return writeString(filename, content, -1)
}

func OverwriteBytes(filename string, content []byte) error {
	return writeString(filename, string(content), -1)
}

func Mkdir(filename string) error {
	if !IsFileExist(filename) {
		err := os.MkdirAll(filepath.Dir(filename), 0777)
		if err != nil {
			println("mk dir failed ", filename, " failed,", err)
			return err
		}
	}
	return nil
}

func Mkfile(filename string) error {
	if IsFileExist(filename) {
		return nil
	}
	err := Mkdir(filename)
	if err != nil {
		return err
	}
	_, err = os.Create(filename)
	return err
}

func writeString(filename string, content string, mode int) error {
	var f *os.File
	var err error
	if err = Mkdir(filename); err != nil {
		return err
	}
	if mode == 1 {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	} else if mode == -1 {
		f, err = os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
	}
	if err != nil {
		println("open file failed ", filename, " failed,", err)
		return err
	}

	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		println("append file failed!", err.Error())
		return err
	}

	return nil
}

func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func DeleteFile(filename string) bool {
	err := os.RemoveAll(filename)
	if err != nil {
		return false
	}
	return true
}

func DeleteFiles(array []string) error {
	for _, path := range array {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func Rename(root, old, new string) error {
	root = getPath(root)
	// println(old, new, path)
	return os.Rename(root+old, root+new)
}

func MV(array []string, dir string) error {
	args := []string{"-f"}
	args = append(args, array...)
	args = append(args, dir)
	// println("->", dir)
	cmd := exec.Command("mv", args...)
	_, err := cmd.CombinedOutput()
	// println(string(b))
	return err
}
