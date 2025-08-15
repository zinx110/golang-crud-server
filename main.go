package main

import (
	"fmt"
	"log"
	"net/http"
	"randomGoServer/db"
	"randomGoServer/users"
)

func main() {

	// initialize db
	db.InitDb()
	defer db.DBpool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/users", users.UsersHandler)
	mux.HandleFunc("/users/", users.UsersHandler)
	mux.HandleFunc("/users/{id}", users.UsersHandler)

	fmt.Println("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
