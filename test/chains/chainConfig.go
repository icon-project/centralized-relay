package chains

import "github.com/icon-project/centralized-relay/test/interchaintest/ibc"

type ChainConfig struct {
	Type           string      `mapstructure:"type"`
	Name           string      `mapstructure:"name"`
	ChainID        string      `mapstructure:"chain_id"`
	Images         DockerImage `mapstructure:"image"`
	Bin            string      `mapstructure:"bin"`
	Bech32Prefix   string      `mapstructure:"bech32_prefix"`
	Denom          string      `mapstructure:"denom"`
	SkipGenTx      bool        `mapstructure:"skip_gen_tx"`
	CoinType       string      `mapstructure:"coin_type"`
	GasPrices      string      `mapstructure:"gas_prices"`
	GasAdjustment  float64     `mapstructure:"gas_adjustment"`
	TrustingPeriod string      `mapstructure:"trusting_period"`
	NoHostMount    bool        `mapstructure:"no_host_mount"`
	BlockInterval  int         `mapstructure:"block_interval"`
}

type DockerImage struct {
	Repository string `mapstructure:"repository"`
	Version    string `mapstructure:"version"`
	UidGid     string `mapstructure:"uid_gid"`
}

func (c *ChainConfig) GetIBCChainConfig(chain *Chain) ibc.ChainConfig {

	return ibc.ChainConfig{
		Type:    c.Type,
		Name:    c.Name,
		ChainID: c.ChainID,
		Images: []ibc.DockerImage{{
			Repository: c.Images.Repository,
			Version:    c.Images.Version,
			UidGid:     c.Images.UidGid,
		}},
		Bin:            c.Bin,
		Bech32Prefix:   c.Bech32Prefix,
		Denom:          c.Denom,
		CoinType:       c.CoinType,
		SkipGenTx:      c.SkipGenTx,
		GasPrices:      c.GasPrices,
		GasAdjustment:  c.GasAdjustment,
		TrustingPeriod: c.TrustingPeriod,
		NoHostMount:    c.NoHostMount,
	}
}
