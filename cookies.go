package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func alreadyLoggedIn(req *http.Request) bool {
	cook, err := req.Cookie("session")
	if err != nil {
		return false
	}
	un := dbSessions[cook.Value]
	_, ok := dbUsers[un]
	return ok
}
func CreateCooke(w http.ResponseWriter) *http.Cookie {
	uId := uuid.NewV4()
	c := &http.Cookie{Name: "session", Value: uId.String()}
	http.SetCookie(w, c)
	return c
}
