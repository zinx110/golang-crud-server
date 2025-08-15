package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"randomGoServer/db"
	"strconv"
	"sync"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var serial = 0
var users = make([]User, 0)
var mu = sync.RWMutex{}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UsersHandler called")

	switch r.Method {
	case http.MethodGet:
		if r.PathValue("id") != "" {
			getUser(w, r)
			return
		}
		getAllUsers(w, r)
	case http.MethodPost:
		addUser(w, r)

	case http.MethodDelete:
		deleteUser(w, r)
	case http.MethodPut, http.MethodPatch:
		updateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getAllUsers called")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rows, err := db.DBpool.Query(r.Context(), "SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var output []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		output = append(output, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("adding user")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	decoder := json.NewDecoder(r.Body)
	var data User
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	mu.Lock()
	serial++
	id := serial
	u := User{
		Id:   id,
		Name: data.Name,
	}
	users = append(users, u)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u)
	fmt.Println("User added successfully:", u)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting user")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	idStr := r.PathValue(("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "Id is invalid", http.StatusBadRequest)
		return
	}
	var selectedUser *User
	err = db.DBpool.QueryRow(r.Context(), "Select * from users where id=$1", id).Scan(&selectedUser.Id, &selectedUser.Name)

	if selectedUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(selectedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("User retrieved successfully:", selectedUser)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("removing user")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "id is invalid", http.StatusBadRequest)
		return
	}
	index := -1
	mu.Lock()

	for i, user := range users {
		if user.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	users = append(users[:index], users[index+1:]...)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updating user")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "id is invalid", http.StatusBadRequest)
		return
	}
	var data User
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	if data.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	index := -1
	mu.Lock()
	for i, user := range users {
		if user.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Update the user's fields
	var selectedUser *User

	users[index].Name = data.Name
	u := users[index]
	selectedUser = &u

	mu.Unlock()

	_ = json.NewEncoder(w).Encode(selectedUser)
}
