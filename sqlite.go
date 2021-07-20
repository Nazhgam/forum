package main

import (
	"database/sql"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func createDatabase() {
	_, err := sql.Open("sqlite3", "./db/datas.db")
	if err != nil {
		os.Create("./db/sqlite-database.db")

	}
	db, _ = sql.Open("sqlite3", "./db/datas.db")
}
func createTables() {
	createDatabase()
	stuser, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS user(id INTEGER PRIMARY KEY, username TEXT, email TEXT UNIQUE, password TEXT, first_name TEXT, last_name TEXT)`)
	stpost, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS post(id INTEGER PRIMARY KEY, user_id INTEGER, categ_id INTEGER, title TEXT, post TEXT,  FOREIGN KEY(user_id) REFERENCES user(id), FOREIGN KEY(categ_id) REFERENCES categ(id))`)
	stcateg, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS categ(id INTEGER PRIMARY KEY, category TEXT UNIQUE)`)
	stcomment, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS comment(id INTEGER PRIMARY KEY, user_id INTEGER, post_id INTEGER, commentary TEXT,  FOREIGN KEY(user_id) REFERENCES user(id), FOREIGN KEY(post_id) REFERENCES post(id))`)
	stlike, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS likes(id INTEGER PRIMARY KEY,user_id INTEGER, post_id INTEGER, FOREIGN KEY(user_id) REFERENCES user(id), FOREIGN KEY(post_id) REFERENCES post(id))`)
	stdislike, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS dislikes(id INTEGER PRIMARY KEY,user_id INTEGER, post_id INTEGER, FOREIGN KEY(user_id) REFERENCES user(id), FOREIGN KEY(post_id) REFERENCES post(id))`)
	stuser.Exec()
	stpost.Exec()
	stcateg.Exec()
	stcomment.Exec()
	stlike.Exec()
	stdislike.Exec()
	insertToCateg()
}

////////////////////////////////////////////////////////////////////////////////////////  INSERT STRUCT INTO THE DATABASE TABLES
func insertToCateg() {
	arr := []string{"food", "music", "movie", "anime", "sport", "math", "golang", "games"}
	for _, i := range arr {
		db.Exec(`INSERT INTO categ (category) VALUES($1)`, i)
	}
}
func insertToUser(u user) {

	db.Exec(`INSERT INTO user (username, email, password, first_name, last_name) VALUES($1, $2, $3, $4, $5)`, u.UserName, u.Email, string(u.Password), u.First, u.Last)

}
func insertToPost(p *post) {

	db.Exec(`INSERT INTO post (user_id, categ_id, title, post) VALUES($1, $2, $3, $4)`, p.UserId, p.CategId, p.Title, p.Text)
}
func insertToComment(c comments) {
	db.Exec(`INSERT INTO commect (user_id, post_id, commentary) VALUES($1, $2, $3)`, c.UserId, c.PostId, c.Comment)
}
func insertToLike(p post) {
	db.Exec(`INSERT INTO likes (user_id, post_id) VALUES($1, $2)`, p.UserId, p.CategId)
}
func insertToDislike(p post) {
	db.Exec(`INSERT INTO dislikes (user_id, post_id) VALUES(1$, $2)`, p.UserId, p.CategId)
}

////////////////////////////////////////////////////////////////////////          PARSE DATAS INTO THE STRUCT
func parseFromCateg() []categ {
	var c []categ
	var cc categ
	row, err := db.Query(`SELECT * FROM categ`)
	if err != nil {
		return nil
	}
	for row.Next() {
		row.Scan(&cc.Id, &cc.Categ)
		c = append(c, cc)
	}
	return c
}
func parseFromUser() {

}
func parseSingleUser(u user) user {
	row := db.QueryRow(`SELECT id, username, first_name, last_name FROM user WHERE email=$1`, u.Email)
	row.Scan(&u.Id, &u.UserName, &u.First, &u.Last)
	return u
}
func parseSinglePost(id int) post {
	var p post
	row := db.QueryRow(`SELECT * FROM post WHERE id = $1`, id)
	row.Scan(&p.Id, &p.UserId, &p.CategId, &p.Title, &p.Text)
	rw, err := db.Query(`SELECT COUNT(id) FROM likes WHERE post_id = $1`, id)
	if err == nil {
		rw.Scan(&p.Like)
	}
	rwo, err := db.Query(`SELECT COUNT(id) FROM dislikes WHERE post_id = $1`, id)
	if err == nil {
		rwo.Scan(&p.DisLike)
	}
	return p
}
func parseFromPost(id int) []post {
	var postes []post
	var pp post
	row, err := db.Query(`SELECT id, user_id, categ_id, title, post FROM post WHERE categ_id=$1`, id)
	if err != nil {
		return nil
	}
	for row.Next() {
		row.Scan(&pp.Id, &pp.UserId, &pp.CategId, &pp.Title, &pp.Text)
		postes = append(postes, pp)
	}
	return postes
}
func parseFromComment() []comments {
	c := comments{PostId: idPost}
	cSlice := []comments{}
	row, _ := db.Query(`SELECT post.id, comment.commentary FROM comment LEFT JOIN post ON comment.post_id = $1`, idPost)
	rwo, _ := db.Query(`SELECT user.username from comment LEFT JOIN user ON comment.user_id = user.id`, idUser)
	for row.Next() && rwo.Next() {
		row.Scan(&c.PostId, &c.Comment)
		rwo.Scan(&c.User)
		cSlice = append(cSlice, c)
	}
	return cSlice
}
func parseFromLike() {

}
func parseFromDislike() {

}

///////////////////////////////////////////////////////////////////////// check user (liked/dislike) or not
func checkUserForLike(id int) bool {
	_, err := db.Query(`SELECT user_id FROM likes WHERE user_id=$1`, id)
	return err == nil
}
func checkUserForDislike(id int) bool {
	_, err := db.Query(`SELECT user_id FROM dislikes WHERE user_id=$1`, id)
	return err == nil
}
func alreadyRegistredInFromDb(u user) bool {
	us := &user{}
	p := ""
	sqlStatement := `SELECT  email, password FROM user WHERE email=?`
	row := db.QueryRow(sqlStatement, u.Email)
	err := row.Scan(&us.Email, &p)
	us.Password = []byte(p)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(us.Password, u.Password)
	return err == nil
}
