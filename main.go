package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Listening on port : 3000")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(*r)
		fmt.Fprintf(w, "<h1>Hello world!</h1>")
	})

	http.ListenAndServe(":3000", nil)
}
