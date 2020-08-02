package api

import (
	"fmt"
	"net/http"
)

const (
	USERNAME = "username"
	VIEW     = "view"
	ID       = "id"
)

var SESSION_MAP = map[string]Session{}

type Session map[string]string

func (s Session) GetUsername() string {
	return s[USERNAME]
}

func (s Session) SetView(view string) {
	s[VIEW] = view
}

func (s Session) GetView() string {
	return s[VIEW]
}

func (s Session) SetId(uuid string) {
	s[ID] = uuid
}

func (s Session) GetId() string {
	return s[ID]
}

func (s Session) Update() {
	SESSION_MAP[s.GetId()] = s
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("SESSION")
	if err == nil && session != nil {
		delete(SESSION_MAP, session.Value)
	}
	//init a new session,(not login)
	w.Header().Set("Set-Cookie", fmt.Sprintf("SESSION=%s", GetUUID()))
}

func GetSession(r *http.Request) Session {
	session, err := r.Cookie("SESSION")
	if err != nil || session == nil {
		return nil
	}
	return SESSION_MAP[session.Value]
}

func GetUsername(r *http.Request) string {

	session := GetSession(r)
	if session == nil || session.GetUsername() == "" {
		return ""
	}
	return session.GetUsername()
}

func SetSession(session Session, w http.ResponseWriter) {
	uuid := GetUUID()
	session.SetId(uuid)
	SESSION_MAP[uuid] = session
	// b,_:=json.Marshal(SESSION_MAP)
	// if !checkFileIsExist("session"){
	// 	os.Create("session")
	// }
	w.Header().Set("Set-Cookie", fmt.Sprintf("SESSION=%s", uuid))
}

func HasSession(r *http.Request) bool {
	return GetUsername(r) != ""
}

func NewSession(username string, view string) Session {
	return map[string]string{USERNAME: username, VIEW: view}
}
