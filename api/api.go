package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eineder/kivi/store"
)

type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

var s store.Store = store.NewInMemoryStore()

func Start(port string) {
	http.HandleFunc("/items/", itemsHandler)

	fmt.Printf("Starting server on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/items/")

	if key == "" {
		http.Error(w, "Key expected", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getItem(w, r, key)
	case http.MethodPost:
		createItem(w, r, key)
	case http.MethodPut:
		updateItem(w, r, key)
	case http.MethodDelete:
		deleteItem(w, r, key)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createItem(w http.ResponseWriter, r *http.Request, key string) {
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = s.CreateItem(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeResponse(w, true, fmt.Sprintf("Item %s created", key), nil)
}

func getItem(w http.ResponseWriter, r *http.Request, key string) {
	value, err := s.GetItem(key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		writeResponse(w, false, err.Error(), nil)
		return
	}
	writeResponse(w, true, "Value for key found", value)
}

func updateItem(w http.ResponseWriter, r *http.Request, key string) {
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = s.UpdateItem(key, string(value))

	if err != nil {
		writeResponse(w, false, fmt.Sprintf("Item %s not found", key), nil)
		return
	}

	writeResponse(w, true, fmt.Sprintf("Item %s updated", key), nil)
}

func deleteItem(w http.ResponseWriter, r *http.Request, key string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := s.DeleteItem(key)

	if err != nil {
		writeResponse(w, false, fmt.Sprintf("Item %s not found", key), nil)
		return
	}
	writeResponse(w, true, fmt.Sprintf("Item %s deleted", key), nil)
}

func writeResponse(w http.ResponseWriter, success bool, message string, payload interface{}) {
	json.NewEncoder(w).Encode(Response{Success: success, Message: message, Payload: payload})
}
