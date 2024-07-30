package bitcoin

import (
	"log"
	"net/http"
	"os"
)

func startSlave() {
	http.HandleFunc("/", handleRoot)
	server := &http.Server{
		Addr:    "8080",
		Handler: nil,
	}

	log.Printf("Slave starting on port %s", "8080")
	log.Fatal(server.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			http.Error(w, "Missing API Key", http.StatusUnauthorized)
			return
		}
		apiKeyHeader := os.Getenv("API_KEY")
		if apiKey != apiKeyHeader {
			http.Error(w, "Invalid API Key", http.StatusForbidden)
			return
		}
		w.Write([]byte("hello world"))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
