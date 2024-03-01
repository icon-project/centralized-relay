package ibc

type RelayerConfig struct {
	Global struct {
		ApiListenAddr  int    `yaml:"api-listen-addr"`
		Timeout        string `yaml:"timeout"`
		Memo           string `yaml:"memo"`
		LightCacheSize int    `yaml:"light-cache-size"`
		KMSKeyID       string `yaml:"kms-key-id"`
	} `yaml:"global"`
	Chains map[string]interface{} `yaml:"chains"`
}

// ChainConfig defines the chain parameters requires to run an interchaintest testnet for a chain.
type ChainConfig struct {
	// Chain type, e.g. cosmos.
	Type string `yaml:"type"`
	// Chain name, e.g. cosmoshub.
	Name string `yaml:"name"`
	// Chain ID, e.g. cosmoshub-4
	ChainID string `yaml:"chain-id"`
	// Docker images required for running chain nodes.
	Images []DockerImage `yaml:"images"`
	// Binary to execute for the chain node daemon.
	Bin string `yaml:"bin"`
	// Bech32 prefix for chain addresses, e.g. cosmos.
	Bech32Prefix string `yaml:"bech32-prefix"`
	// Denomination of native currency, e.g. uatom.
	Denom string `yaml:"denom"`
	// Coin type
	CoinType string `default:"118" yaml:"coin-type"`
	// Minimum gas prices for sending transactions, in native currency denom.
	GasPrices string `yaml:"gas-prices"`
	// Adjustment multiplier for gas fees.
	GasAdjustment float64 `yaml:"gas-adjustment"`
	// Trusting period of the chain.
	TrustingPeriod string `yaml:"trusting-period"`
	// Do not use docker host mount.
	NoHostMount bool `yaml:"no-host-mount"`
	// When true, will skip validator gentx flow
	SkipGenTx bool
	// When provided, will run before performing gentx and genesis file creation steps for validators.
	PreGenesis func(ChainConfig) error
	// When provided, genesis file contents will be altered before sharing for genesis.
	ModifyGenesis func(ChainConfig, []byte) ([]byte, error)
	// Override config parameters for files at filepath.
	ConfigFileOverrides map[string]any
}

type DockerImage struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
	UidGid     string `yaml:"uid-gid"`
}

// Ref returns the reference to use when e.g. creating a container.
func (i DockerImage) Ref() string {
	if i.Version == "" {
		return i.Repository + ":latest"
	}

	return i.Repository + ":" + i.Version
}

type WalletAmount struct {
	Address string
	Denom   string
	Amount  int64
}

type Wallet interface {
	KeyName() string
	FormattedAddress() string
	Mnemonic() string
	Address() []byte
}
