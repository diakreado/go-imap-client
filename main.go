package main

import (
	"fmt"
	"html/template"
	"net/http"

	db "./db"
	imap "./imap"
	. "github.com/logrusorgru/aurora"
)

const (
	port = ":3000"
)

// Index(roote) route handler
// read auth data from auth.json and show login/logout info
// also view result of requset to imap server
func indexHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		"./templates/index.html",
		"./templates/header.html",
		"./templates/footer.html",
		"./templates/content.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println(req.Method, req.URL)

	data := db.GetAuthData()

	if data.Login != "" && data.Password != "" && data.Server != "" {
		imap.GetPostBoxState()
		// fmt.Println("Succsessful authentication!")
	}

	t.ExecuteTemplate(w, "index", data)
}

// Auth route hanlder
// write data from Post request to auth.json
// and send redirect to index(root)
func authHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	data := db.AuthData{
		Login:    req.FormValue("login"),
		Password: req.FormValue("password"),
		Server:   req.FormValue("server")}

	db.SaveAuthData(data)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)
	fmt.Println("Listening on port", Brown(port))
	err := http.ListenAndServe(port, nil)
	fmt.Println("Error creating http server:", Red(err))
}
