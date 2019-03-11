package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// TODO
	// Прочитать данные и отправить их в форму для рендеринга, если там пусто то просил ввести их,
	// если нет, то рисуем данные
	// данные передаём вместо nil
	t.ExecuteTemplate(w, "index", nil)
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	form := req.Form

	fo, err := os.Create("auth")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	enc := json.NewEncoder(fo)
	enc.Encode(form)
}

func main() {
	fmt.Println("Listening on port :3000")

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)
	http.ListenAndServe(":3000", nil)
}
