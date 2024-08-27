package bitcoin

import (
	"io"
	"log"
	"net/http"
	"os"
)

func startSlave(c *Config) {
	http.HandleFunc("/", handleRoot)
	port := c.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	log.Printf("Slave starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		log.Printf("Received payload: %s", string(body))
		w.Write([]byte("Payload received" + string(body)))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
