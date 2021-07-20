package main

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, req *http.Request) {
	// CreateCooke(w, req)
	c := parseFromCateg()
	if req.Method == http.MethodPost {
		idCateg, _ = strconv.Atoi(req.FormValue("Id"))
		http.Redirect(w, req, "/theme", http.StatusSeeOther)
	}
	tmp.ExecuteTemplate(w, "index.html", c)
}

///////////////////////////////
func theme(w http.ResponseWriter, req *http.Request) {
	p := parseFromPost(idCateg)
	if req.Method == http.MethodPost {
		idPost, _ = strconv.Atoi(req.FormValue("id"))
		http.Redirect(w, req, "/postes", http.StatusSeeOther)
		return
	}
	tmp.ExecuteTemplate(w, "theme.html", p)
}

/////////////////////////
func postes(w http.ResponseWriter, req *http.Request) {
	p := parseSinglePost(idPost)
	if req.Method == http.MethodPost {
		like, _ := strconv.Atoi(req.FormValue("postIdCom"))
		dis, _ := strconv.Atoi(req.FormValue("postIdCom"))
		idPost, _ = strconv.Atoi(req.FormValue("postIdCom"))
		if !checkUserForLike(like) {
			insertToLike(p)
		}
		if !checkUserForDislike(dis) {
			insertToDislike(p)
		}
	}
	tmp.ExecuteTemplate(w, "postes.html", p)
}

//////////////////////////////////
func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, err := req.Cookie("session")
	if err != nil {
		c = CreateCooke(w)
	}
	if req.Method == http.MethodPost {
		un := req.FormValue("nickname")
		email := req.FormValue("username")
		p := req.FormValue("password")
		fn := req.FormValue("firstname")
		ln := req.FormValue("lastname")
		if !isValidPassword(p) {
			http.Redirect(w, req, "/signup", http.StatusSeeOther)
			return
		}
		bp, err := bcrypt.GenerateFromPassword([]byte(p), 5)
		if err != nil {
			return
		}

		u := user{UserName: un, Email: email, First: fn, Last: ln, Password: bp}
		dbSessions[c.Value] = un
		dbUsers[un] = u
		if alreadyRegistredInFromDb(u) {
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		insertToUser(u)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tmp.ExecuteTemplate(w, "signup.html", nil)
}

///////////////////// check password, valid or not
func isValidPassword(p string) bool {
	if len(p) < 5 {

		return false
	}
	if !strings.ContainsAny(p, "ABCDEFGHIJKLMNOPQRSTUVWXYZ,.!@#$%^&*(){}_+=-0123456789abcdefghijklmnopqrstuvwxyz") {

		return false
	}
	return true
}

/////////////////////////
func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, err := req.Cookie("session")
	if err != nil {
		c = CreateCooke(w)
	}
	if req.Method == http.MethodPost {
		email := req.FormValue("username")
		p := req.FormValue("password")
		u := user{Email: email, Password: []byte(p)}
		if !alreadyRegistredInFromDb(u) {
			http.Redirect(w, req, "/signup", http.StatusSeeOther)
			return
		}
		u = parseSingleUser(u)
		idUser = u.Id
		dbSessions[c.Value] = u.UserName
		dbUsers[u.UserName] = u
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return

	}
	tmp.ExecuteTemplate(w, "login.html", nil)
}

/////////////////////////
func comment(w http.ResponseWriter, req *http.Request) {
	arrCom := parseFromComment()
	if req.Method == http.MethodPost {

	}
	tmp.ExecuteTemplate(w, "comment.html", arrCom)
}

/////
func title(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	cc := parseFromCateg()
	if req.Method == http.MethodPost {
		title := req.FormValue("title")
		text := req.FormValue("text")
		idcateg, _ := strconv.Atoi(req.FormValue("id"))
		p := &post{UserId: idUser, CategId: idcateg, Title: title, Text: text}
		insertToPost(p)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tmp.ExecuteTemplate(w, "title.html", cc)
}
