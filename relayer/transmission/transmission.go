package transmission

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func CallBitcoinRelay(message string) []byte {
	fmt.Printf("Call Bitcoin Relayer")
	masterServer := os.Getenv("MASTER_SERVER")

	client := &http.Client{}
	apiKeyHeader := os.Getenv("API_KEY")

	req, err := http.NewRequest("POST", masterServer+"/execute", bytes.NewBuffer([]byte(message)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("x-api-key", apiKeyHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	return body
}
