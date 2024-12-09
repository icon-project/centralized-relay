package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func startMaster(c *Config) {
	http.HandleFunc("/execute", handleExecute)
	port := c.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	log.Printf("Master starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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

	var msg string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "msg": msg}
	json.NewEncoder(w).Encode(response)
}

func requestPartialSign(apiKey string, url string, slaveRequestData []byte, responses chan<- slaveResponse, order int, wg *sync.WaitGroup) {
	defer wg.Done()
	response := slaveResponse{}
	client := &http.Client{}
	payload := bytes.NewBuffer(slaveRequestData)
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		response.err = fmt.Errorf("failed to create request: %v", err)
		responses <- response
		return
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)

	if err != nil {
		response.err = fmt.Errorf("failed to send request: %v", err)
		responses <- response
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.err = fmt.Errorf("error reading response: %v", err)
		responses <- response
		return
	}

	sigs := [][]byte{}
	err = json.Unmarshal(body, &sigs)
	if err != nil {
		response.err = fmt.Errorf("err Unmarshal: %v", err)
		responses <- response
		return
	}

	response.order = order
	response.sigs = sigs
	response.err = nil
	responses <- response
}
