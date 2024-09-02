package bitcoin

import (
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
)

func startSlave(c *Config, p *Provider) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRoot(w, r, p)
	})
	port := c.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	log.Printf("Slave starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request, p *Provider) {
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

		var rsi slaveRequestParams
		err = json.Unmarshal(body, &rsi)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			return
		}
		sigs, _ := loadSigsFromDb(rsi.MsgSn, p)
		// return sigs to master
		returnData, _ := json.Marshal(sigs)
		w.Write(returnData)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func loadSigsFromDb(sn *big.Int, p *Provider) ([][]byte, error) {
	key := sn.String()
	data, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var sigs [][]byte
	err = json.Unmarshal(data, &sigs)
	if err != nil {
		return nil, err
	}
	return sigs, nil
}
