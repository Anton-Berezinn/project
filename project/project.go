package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

var LogFile = "log.txt"

type Log struct {
	log *log.Logger
}

func (log *Log) MainPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.New("").ParseFiles("./template/main.html")
	if err != nil {
		log.log.Println("err in open login.html\nErr:", err)
	}
	err = tmp.ExecuteTemplate(w, "main.html", nil)
	if err != nil {
		log.log.Println("err in execute login.html\nErr:", err)
	}

}

type User struct {
	Id       int
	Username string
	Password string
}

func (log *Log) ErrorPasword(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Wrong Password</h1>")
}

func (log *Log) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.New("").ParseFiles("./template/login.html")
	if err != nil {
		log.log.Println("err in open login.html\nErr:", err)
		return
	}
	err = tmp.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.log.Println("err in execute login.html\nErr:", err)
		return
	}
	connstr := "user=" + os.Getenv("Username") + " password=" + os.Getenv("password") + " dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	rows, err := db.Query("SELECT * FROM users WHERE Username = $1", username)
	if err != nil {
		log.log.Println("err in prepare statement\nErr:", err)
		return
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.log.Println("err in scan\nErr:", err)
			return
		}
		// теперь user содержит данные из базы данных
	}
	if password != user.Password {
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}
}

func main() {

	f, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	ilog := log.New(f, "customLogLineNumber", log.LstdFlags)
	ilog.SetFlags(log.Lshortfile)
	defer f.Close()
	er := &Log{
		log: ilog,
	}

	http.HandleFunc("/", er.MainPage)
	http.HandleFunc("/login", er.LoginPage)
	http.HandleFunc("/error", er.ErrorPasword)
	fmt.Println(http.ListenAndServe("localhost:8080", nil), "start serve")
}
