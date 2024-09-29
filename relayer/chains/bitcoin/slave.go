package bitcoin

import (
	"encoding/json"
	"io"
	"net/http"

	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
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

	p.logger.Info("Slave starting on port", zap.String("port", port))
	p.logger.Fatal("Failed to start slave", zap.Error(server.ListenAndServe()))
}

func handleRoot(w http.ResponseWriter, r *http.Request, p *Provider) {
	p.logger.Info("Slave starting on port", zap.String("port", p.cfg.Port))
	if r.Method == http.MethodPost {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			p.logger.Error("Missing API Key")
			http.Error(w, "Missing API Key", http.StatusUnauthorized)
			return
		}
		apiKeyHeader := p.cfg.ApiKey
		if apiKey != apiKeyHeader {
			p.logger.Error("Invalid API Key", zap.String("apiKey", apiKey))
			http.Error(w, "Invalid API Key", http.StatusForbidden)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			p.logger.Error("Error reading request body", zap.Error(err))
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var rsi slaveRequestParams
		err = json.Unmarshal(body, &rsi)
		if err != nil {
			p.logger.Error("Error decoding request body", zap.Error(err))
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			return
		}
		sigs, _ := buildAndSignTxFromDbMessage(rsi.MsgSn, p)
		// return sigs to master
		returnData, _ := json.Marshal(sigs)
		w.Write(returnData)
	} else {
		p.logger.Error("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func buildAndSignTxFromDbMessage(sn string, p *Provider) ([][]byte, error) {
	key := sn
	data, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var message *relayTypes.Message
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}

	_, _, _, relayerSigns, _, _, err := p.HandleBitcoinMessageTx(message)
	if err != nil {
		return nil, err
	}

	return relayerSigns, nil
}
