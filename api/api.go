package api

import (
	"fmt"
	"net/http"
	"strings"
)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"message": "Item %s created"}`, key)))
}

func getItem(w http.ResponseWriter, r *http.Request, key string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"item": {"key": "%s"}}`, key)))
}

func updateItem(w http.ResponseWriter, r *http.Request, key string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Item %s updated"}`, key)))
}

func deleteItem(w http.ResponseWriter, r *http.Request, key string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Item %s deleted"}`, key)))
}
