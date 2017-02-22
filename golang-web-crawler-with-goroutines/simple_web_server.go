// based on: http://stackoverflow.com/questions/42381426/change-the-sample-by-using-goroutine
package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("Starting server in port: 8080")
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is about us: /about")
	}
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to: /home")
	}


	http.HandleFunc("/about", rootHandler)
	http.HandleFunc("/home", homeHandler)
	http.ListenAndServe(":8080", nil)
}
