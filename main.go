package main

import (
	"fmt"
	"log"
	"net/http"
	"randomGoServer/users"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/users", users.UsersHandler)
	mux.HandleFunc("/users/", users.UsersHandler)
	mux.HandleFunc("/users/{id}", users.UsersHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
