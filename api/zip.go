package api

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var TMP_ROOT = ".tmp/"

func setTMP_ROOT(dir string) {
	TMP_ROOT = getPath(dir)
}

func Zip(root string, maxFilter int64, fileList []string) (string, error) {
	totalSize := int64(0)
	for _, filePath := range fileList {
		totalSize += FileSize(filePath)
	}
	if totalSize > maxFilter {
		return "", errors.New("文件总大小超过" + formatSize(totalSize))
	}
	//开始压缩
	// 截取文件列表的前10个
	totalName := strings.Join(fileList, ",")
	tmpName := GetMd5String(totalName) + ".zip" //生成文件名的md5
	zipPath := root + TMP_ROOT
	os.MkdirAll(zipPath, 0777) //创建目录
	zipFileName := zipPath + tmpName
	// 如果存在则不创建
	if IsFileExist(zipFileName) {
		return tmpName, nil
	}
	zipfile, err := os.Create(zipFileName)
	if err != nil {
		return "", err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	for _, srcFile := range fileList {

		filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
			// header.Name = path
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
			}
			return err
		})
	}
	return tmpName, err
}
