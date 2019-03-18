package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	db "./db"
	imap "./imap"
	. "github.com/logrusorgru/aurora"
)

const (
	port = ":3000"
)

type Data struct {
	Auth     db.AuthData
	Envelope []imap.Envelope
}

var funcMap = template.FuncMap{
	"trunc": func(c int, s string) string {
		runes := []rune(s)
		if len(runes) <= c {
			return s
		}
		return string(runes[:c]) + "..."
	},
	"dec": func(i int) int {
		return i - 1
	},
}

func isLocalhost(remoteAddr string) bool {
	var localhost = regexp.MustCompile(`127.0.0.1:`)
	return localhost.MatchString(remoteAddr)
	// return true
}

// Index(roote) route handler
// read auth data from auth.json and show login/logout info
// also view result of requset to imap server
func indexHandler(res http.ResponseWriter, req *http.Request) {
	if !isLocalhost(req.RemoteAddr) {
		http.Redirect(res, req, "https://http://hyperborea-theatre.ru/", http.StatusSeeOther)
		return
	}
	fmt.Println(req.Method, req.URL)

	templates := template.Must(template.New("main").Funcs(funcMap).ParseGlob("./templates/*"))

	authData := db.GetAuthData()

	var data Data
	if authData.Login != "" && authData.Password != "" && authData.Server != "" {
		envelopeData := imap.GetListOfMails()

		data.Auth = authData
		data.Envelope = envelopeData
	}

	templates.ExecuteTemplate(res, "index", data)
}

func letterHandler(res http.ResponseWriter, req *http.Request) {
	if !isLocalhost(req.RemoteAddr) {
		http.Redirect(res, req, "https://http://hyperborea-theatre.ru/", http.StatusSeeOther)
		return
	}
	uid := req.FormValue("uid")
	fmt.Println(req.Method, req.URL)

	templates := template.Must(template.New("main").Funcs(funcMap).ParseGlob("./templates/*"))

	letter := imap.GetLetter(uid)

	templates.ExecuteTemplate(res, "letter", letter)
}

// Auth route hanlder
// write data from Post request to auth.json
// and send redirect to index(root)
func authHandler(res http.ResponseWriter, req *http.Request) {
	if !isLocalhost(req.RemoteAddr) {
		http.Redirect(res, req, "https://http://hyperborea-theatre.ru/", http.StatusSeeOther)
		return
	}
	fmt.Println(req.Method, req.URL)
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	data := db.AuthData{
		Login:    req.FormValue("login"),
		Password: req.FormValue("password"),
		Server:   req.FormValue("server")}

	db.SaveAuthData(data)

	if req.Method == "POST" && data.Login != "" && data.Password != "" && data.Server != "" {
		fmt.Fprintf(res, "%t", imap.TryToLogin())
	} else {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}

func faviconHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "lol")
}

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/letter", letterHandler)

	fmt.Println("Listening on port", Brown(port))
	err := http.ListenAndServe(port, nil)
	fmt.Println("Error creating http server:", Red(err))
}
