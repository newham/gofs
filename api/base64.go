package api

import (
	"encoding/base64"
)

func URLToBase64(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func Base64ToURL(str string) string {
	b, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(b)
}
