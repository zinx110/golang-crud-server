package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Id   int    `json:"id,string"`
	Name string `json:"name"`
}

var users = make([]User, 0)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
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
	encoder := json.NewEncoder(w)
	err := encoder.Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("adding user")
	decoder := json.NewDecoder(r.Body)
	var data User
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.Id = len(users) + 1
	if data.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	users = append(users, data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(data)
	fmt.Println("User added successfully:", data)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting user")
	idStr := r.PathValue(("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "Id is invalid", http.StatusBadRequest)
		return
	}
	for _, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
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
	for i, user := range users {
		if user.Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	users = append(users[:index], users[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updating user")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "id is invalid", http.StatusBadRequest)
		return
	}
	index := -1
	for i, user := range users {
		if user.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var data User
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	// Update the user's fields
	if data.Name != "" {
		users[index].Name = data.Name
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(users[index])
}
