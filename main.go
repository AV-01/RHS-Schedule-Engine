package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello, RaspAPI!",
		})
	})

	mux.HandleFunc("GET /greet", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if name == "" {
			name = "world"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Hello, %s", name),
		})
	})

	var things = map[string]string{
		"1": "first thing",
		"2": "Second thing",
	}

	mux.HandleFunc("GET /things/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		thing, ok := things[id]

		if !ok {
			writeError(w, http.StatusNotFound, "thing not found")
			return
		}

		writeJSON(w, http.StatusOK, thing)
	})

	type JellyBeans struct {
		Flavor   string `json:"flavor"`
		Color    string `json:"color"`
		Quantity int    `json:"quantity"`
	}

	mux.HandleFunc("POST /eatbeans", func(w http.ResponseWriter, r *http.Request) {
		var beans JellyBeans

		if err := json.NewDecoder(r.Body).Decode(&beans); err != nil {
			http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("You ate %d %s %s jellybeans!", beans.Quantity, beans.Color, beans.Flavor),
		})

	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs.html")
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}
