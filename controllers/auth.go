package controllers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/newham/gofs/api"
)

var TOKEN_MAP = map[string]string{}
var TOKEN_INDEX_MAP = map[string]string{}

func AuthController(w http.ResponseWriter, r *http.Request) {
	userName := api.GetSession(r).GetUsername()
	token := TOKEN_MAP[userName]
	if token == "" || r.FormValue("f") == "true" {
		token = fmt.Sprintf("%x", sha256.Sum256([]byte(userName+strconv.FormatInt(time.Now().UnixNano(), 10))))
		TOKEN_MAP[userName] = token
		TOKEN_INDEX_MAP[token] = userName
	}
	log(200, "auth", "user:"+userName)
	w.WriteHeader(200)
	w.Write([]byte(token))
}

func checkToken(r *http.Request) bool {
	//3.check it
	if TOKEN_INDEX_MAP[getToken(r)] != "" {
		return true
	}
	return false
}

func getToken(r *http.Request) string {
	//1.by head
	token := r.Header.Get("auth-token")
	//2.by param
	if token == "" {
		token = r.FormValue("auth-token")
	}
	return token
}

func getUsernameByToken(r *http.Request) string {
	return TOKEN_INDEX_MAP[getToken(r)]
}
