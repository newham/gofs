/******************************************
*FileName: user.go
*Author: Liu han
*Date: 2020-08-18
*Description: 简单的基于csv的用户数据库服务
*******************************************/
package api

import (
	"encoding/json"
	"time"
)

var USER_DB, _ = NewCSV("./db/user.csv")

type User struct {
	ID         string
	Name       string
	PWD        string
	Email      string
	Phone      string
	Info       string
	Status     int
	CreateTime int64
	Rights     string
}

func NewUser(Name string,
	PWD string,
	Email string,
	Phone string,
	Info string) User {
	return User{GetMd5String(Name), Name, PWD, Email, Phone, Info, 1, time.Now().Unix(), "w,r"}
}

func (u User) ToJSONString() string {
	b, err := json.Marshal(u)
	if err != nil {
		println(err.Error)
		return ""
	}
	return URLToBase64(string(b))
}

func PutUser(user User) error {
	l := NewLine(user.ID, user.ToJSONString())
	USER_DB.Put(l)
	return USER_DB.Save()
}

func GetUser(id string) *User {
	line := USER_DB.Get(id)
	if line == nil {
		return nil
	}
	jsonStr := Base64ToURL(line[1])
	user := User{}
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		println(err.Error())
		return nil
	}
	return &user
}

//这里采用布隆过滤器进行校验
