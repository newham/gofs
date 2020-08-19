package api

import (
	"testing"
	"time"
)

var uuid = "73fd3bbe5110d029d6e23eedf70efe0a"

func TestPutUser(t *testing.T) {
	PutUser(User{uuid, "liuhan", "123", "test@test.com", "15821123123", "我就是我", 1, time.Now().Unix(), "w,r,x"})
}

func TestGetUser(t *testing.T) {
	user := GetUser(uuid)
	println(user.Email)
}
