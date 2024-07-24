package bitcoin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func startMaster() {
	go requestSlaves()
	masterPort := os.Getenv("MASTER_PORT")
	server := &http.Server{
		Addr:    masterPort,
		Handler: nil,
	}

	log.Printf("Master starting on port %s", masterPort)
	log.Fatal(server.ListenAndServe())
}

func requestSlaves() {
	fmt.Printf("Master request slave")
	slavePort := os.Getenv("SLAVE_PORT")

	// Call slave to get more data
	var wg sync.WaitGroup
	responses := make(chan string, 2)

	wg.Add(1)
	go requestSlave(slavePort, responses, &wg)

	go func() {
		wg.Wait()
		close(responses)
	}()

	for res := range responses {
		fmt.Println("Received response from slave:", res)
	}
}

func requestSlave(url string, responses chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	apiKeyHeader := os.Getenv("API_KEY")

	req, err := http.NewRequest("GET", url, nil)
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
