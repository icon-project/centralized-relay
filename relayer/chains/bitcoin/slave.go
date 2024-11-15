package bitcoin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func startSlave(c *Config, p *Provider) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRoot(w, r, p)
	})
	http.HandleFunc("/update-relayer-message-status", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateRelayerMessageStatus(w, r, p)
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
	if !validateMethod(w, r, p) {
		return
	}
	if !authorizeRequest(w, r, p) {
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
	sigs, _ := buildAndSignTxFromDbMessage(rsi.MsgSn, rsi.FeeRate, p)
	// return sigs to master
	returnData, _ := json.Marshal(sigs)
	w.Write(returnData)
}

func handleUpdateRelayerMessageStatus(w http.ResponseWriter, r *http.Request, p *Provider) {
	if !validateMethod(w, r, p) {
		return
	}
	if !authorizeRequest(w, r, p) {
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.logger.Error("Error reading request body", zap.Error(err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var rsi slaveRequestUpdateRelayMessageStatus
	err = json.Unmarshal(body, &rsi)
	if err != nil {
		p.logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	p.logger.Info("Slave update relayer message status", zap.String("sn", rsi.MsgSn))
	p.updateRelayerMessageStatus(rsi)
}

func authorizeRequest(w http.ResponseWriter, r *http.Request, p *Provider) bool {
	apiKey := r.Header.Get("x-api-key")
	if apiKey == "" {
		p.logger.Error("Missing API Key")
		http.Error(w, "Missing API Key", http.StatusUnauthorized)
		return false
	}
	apiKeyHeader := p.cfg.ApiKey
	if apiKey != apiKeyHeader {
		p.logger.Error("Invalid API Key", zap.String("apiKey", apiKey))
		http.Error(w, "Invalid API Key", http.StatusForbidden)
		return false
	}
	return true
}

func validateMethod(w http.ResponseWriter, r *http.Request, p *Provider) bool {
	if r.Method != http.MethodPost {
		p.logger.Error("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func buildAndSignTxFromDbMessage(sn string, feeRate int, p *Provider) ([][]byte, error) {
	p.logger.Info("Slave start to build and sign tx from db message", zap.String("sn", sn))
	key := sn
	data, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	if strings.Contains(sn, "RB") {
		p.logger.Info("Rollback message", zap.String("sn", sn))
		return nil, nil
	}

	var message *StoredRelayMessage
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}

	_, _, _, relayerSigns, _, err := p.HandleBitcoinMessageTx(message.Message, feeRate)
	if err != nil {
		return nil, err
	}

	return relayerSigns, nil
}

func (p *Provider) updateRelayerMessageStatus(params slaveRequestUpdateRelayMessageStatus) (bool, error) {
	p.logger.Info("Slave update relayer message status", zap.String("sn", params.MsgSn))
	key := params.MsgSn
	data, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return false, err
	}

	if strings.Contains(params.MsgSn, "RB") {
		p.logger.Info("Rollback message", zap.String("sn", params.MsgSn))
		return true, nil
	}

	var message *StoredRelayMessage
	err = json.Unmarshal(data, &message)
	if err != nil {
		return false, err
	}

	message.Status = params.Status
	message.TxHash = params.TxHash

	value, _ := json.Marshal(message)
	err = p.db.Put([]byte(key), value, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
