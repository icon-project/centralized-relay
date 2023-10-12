package icon

import (
	"fmt"

	"go.uber.org/zap"
)

func GetMockIconProvider() (*IconProvider, error) {

	pc := IconProviderConfig{
		ChainID:         "icon",
		KeyStore:        testkeyAddr,
		RPCAddr:         "http://localhost:9082/api/v3",
		Password:        "x",
		StartHeight:     0,
		ContractAddress: "cx0000",
	}
	logger := zap.NewNop()
	prov, err := pc.NewProvider(logger, "", true, "icon-1")
	if err != nil {
		return nil, err
	}

	iconProvider, ok := prov.(*IconProvider)
	if !ok {
		return nil, fmt.Errorf("unbale to type case to icon chain provider")
	}
	return iconProvider, nil

}
