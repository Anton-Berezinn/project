package main

//type User struct {
//	Id   int
//	Name string
//}

//func (u *User) Add() string {
//	if u.Id == 1 {
//
//	}
//	return "hello everyone"
//}
//
//func Useradd(u *User) string {
//	return "hello " + u.Name + "!!"
//}
//
//func mainPage(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("main", "was")
//	//tmpl, err := template.New("").ParseFiles("./template/main.html")
//	//if err != nil {
//	//log.Println("err in parse files", err)
//	//}
//	user := []User{
//		User{1, "anton"},
//		User{2, "nastya"},
//	}
//	tmpl2 := template.FuncMap{
//		"add": Useradd,
//	} //Добавили в шаблон функцию не относящиюся к структуре
//	tmpl, err := template.
//		New("").
//		Funcs(tmpl2).
//		ParseFiles("./template/main.html")
//	if err != nil {
//		log.Println("err in parse", err)
//	}
//	err = tmpl.ExecuteTemplate(w, "main.html", struct {
//		Ann []User
//	}{
//		user,
//	})
//}
//
//func loginPage(w http.ResponseWriter, r *http.Request) {
//	exp := time.Now().Add(30 * time.Minute)
//	cookie := http.Cookie{
//		Name:    "session",
//		Value:   "value",
//		Expires: exp}
//	http.SetCookie(w, &cookie)
//	http.Redirect(w, r, "/", http.StatusFound)
//}
//
//func logoutPage(w http.ResponseWriter, r *http.Request) {
//	cook, err := r.Cookie("session")
//	if err == http.ErrNoCookie {
//		http.Redirect(w, r, "/login", http.StatusFound)
//	} else {
//		cook.Expires = time.Now().AddDate(0, 0, -1)
//		http.SetCookie(w, cook)
//		fmt.Fprintf(w, "<h1>Hello World</h1>"+cook.Name)
//	}
//}
//
//func photoPage(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		src, _, err := r.FormFile("my-file")
//		if err != nil {
//			http.Error(w, err.Error(), 500)
//			return
//		}
//		defer src.Close()
//
//		dst, err := os.Create("file.jpg")
//		if err != nil {
//			http.Error(w, err.Error(), 500)
//			return
//		}
//		defer dst.Close()
//
//		io.Copy(dst, src)
//		fmt.Fprintf(w, "<h1>file was success</h1>")
//	}
//}
//
//func getPhot(w http.ResponseWriter, r *http.Request) {
//	r.ParseMultipartForm(32 << 20)
//	file, handler, err := r.FormFile("my_file")
//	if err != nil {
//		log.Println(err)
//		http.Error(w, "Error parsing file", http.StatusBadRequest)
//		return
//	}
//	defer file.Close()
//	fmt.Println(handler)
//}
//func redirect(w http.ResponseWriter, r *http.Request) {
//	tmpl, err := template.New("").ParseFiles("./template/photo.html")
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//	}
//	err = tmpl.ExecuteTemplate(w, "photo.html", nil)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//	}
//}
//func main() {
//	http.HandleFunc("/lss", redirect)
//	http.HandleFunc("/", mainPage)
//	http.HandleFunc("/upload", photoPage)
//	http.HandleFunc("/get", getPhot)
//	http.HandleFunc("/login", loginPage)
//	http.HandleFunc("/logout", logoutPage)
//	fmt.Println("start work :8080")
//	http.ListenAndServe("localhost:8080", nil)
//}

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func handle(w http.ResponseWriter, req *http.Request) {
	s := ""
	for i := 0; i < 1000; i++ {
		p := &Post{ID: i, Text: "new post"}
		s += fmt.Sprintf("%#v", p)
	}
	w.Write([]byte(s))
}

func main() {
	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
