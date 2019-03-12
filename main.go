package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	db "./db"
)

// Index(roote) route handler
// read auth data from auth.json and show login/logout info
// also view result of requset to imap server
func indexHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./templates/index.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	mess := db.GetAuthData()
	fmt.Println(mess)

	t.ExecuteTemplate(w, "index", mess)
}

// Auth route hanlder
// write data from Post request to auth.json
// and send redirect to index(root)
func authHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	dataForm := map[string]string{
		"login":    req.FormValue("login"),
		"password": req.FormValue("password"),
		"server":   req.FormValue("server")}

	fo, err := os.Create("auth.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	enc := json.NewEncoder(fo)
	enc.Encode(dataForm)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func main() {
	fmt.Println("Listening on port :3000")

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)
	http.ListenAndServe(":3000", nil)
}
