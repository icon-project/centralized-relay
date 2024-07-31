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

func startMaster() {
	http.HandleFunc("/execute", handleExecute)
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	log.Printf("Master starting on port %s", port)

	isProcess := os.Getenv("IS_PROCESS")
	if isProcess == "1" {
		callSlaves("test")
	}
	log.Fatal(server.ListenAndServe())
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// Process the message as needed
	fmt.Println("Received message: ", msg)

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "msg": msg}
	fmt.Println(response)
	json.NewEncoder(w).Encode(response)
}

func callSlaves(txId string) {
	fmt.Println("Master request slave")
	slavePort1 := os.Getenv("SLAVE_SERVER_1")
	slavePort2 := os.Getenv("SLAVE_SERVER_2")
	// Call slave to get more data
	var wg sync.WaitGroup
	responses := make(chan string, 2)

	wg.Add(2)
	go requestPartialSign(slavePort1, txId, responses, &wg)
	go requestPartialSign(slavePort2, txId, responses, &wg)

	go func() {
		wg.Wait()
		close(responses)
	}()

	for res := range responses {
		fmt.Println("Received response from slave: ", res)
	}
}

func requestPartialSign(url string, txId string, responses chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	apiKeyHeader := os.Getenv("API_KEY")
	payload := bytes.NewBuffer([]byte(txId))
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("x-api-key", apiKeyHeader)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	responses <- string(body)
}
