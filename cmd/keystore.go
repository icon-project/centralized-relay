package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/keys"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/spf13/cobra"
)

var TempDir = os.TempDir()

type keystoreState struct {
	chain           string
	password        string
	confirmPassword string
	address         string
	path            string
}

func newKeyStoreState() (*keystoreState, error) {
	return new(keystoreState), nil
}

func keystoreCmd(a *appState) *cobra.Command {
	ks := &cobra.Command{
		Use:     "keystore",
		Aliases: []string{"ks"},
		Short:   "keystore utilty",
		Args:    withUsage(cobra.MaximumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s keystore [command]`, appName)),
	}
	state, err := newKeyStoreState()
	if err != nil {
		panic(err)
	}

	ks.AddCommand(state.init(a), state.new(a), state.list(a), state.importKey(a), state.use(a), state.generateClusterKey(a), state.getClusterKey(a))

	return ks
}

func (k *keystoreState) init(a *appState) *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyID, err := a.kms.Init(cmd.Context())
			if err != nil {
				return err
			}
			a.config.Global.KMSKeyID = *keyID
			if err := a.config.Save(a.configPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "KMS key created: %s\n", a.config.Global.KMSKeyID)
			return nil
		},
	}
	return init
}

// generate ecdsa private key
func (k *keystoreState) generateClusterKey(a *appState) *cobra.Command {
	generate := &cobra.Command{
		Use:   "gen-cluster-key",
		Short: "generate cluster key",
		RunE: func(cmd *cobra.Command, args []string) error {
			keypair, err := keys.NewKeyPair(keys.Secp256k1)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(keys.GetClusterKeyDir(a.homePath), 0o755); err != nil {
				return err
			}

			cipherPk, err := a.kms.Encrypt(context.Background(), keypair.PrivateKey())
			if err != nil {
				return err
			}

			keypath := keys.GetClusterKeyPath(a.homePath, keypair.PublicKey().String())
			if err := os.WriteFile(keypath, cipherPk, 0o600); err != nil {
				return err
			}

			clusterCfg := ClusterConfig{
				Enabled: true,
				PubKey:  keypair.PublicKey().String(),
				keypair: keypair,
			}

			a.config.Global.ClusterConfig = &clusterCfg

			if err := a.config.Save(a.configPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Cluster key created and encrypted: %s\n", a.config.Global.ClusterConfig.PubKey)
			return nil
		},
	}
	return generate
}

// get cluster public key
func (k *keystoreState) getClusterKey(a *appState) *cobra.Command {
	get := &cobra.Command{
		Use:   "get-cluster-key",
		Short: "get cluster key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !a.config.Global.ClusterConfig.Enabled {
				return fmt.Errorf("cluster mode not enabled")
			}
			fmt.Fprintf(os.Stdout, "Cluster key: %s\n", a.config.Global.ClusterConfig.PubKey)
			return nil
		},
	}
	return get
}

func (k *keystoreState) new(a *appState) *cobra.Command {
	new := &cobra.Command{
		Use:   "new",
		Short: "new keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := a.config.Chains.Get(k.chain)
			if err != nil {
				return fmt.Errorf("chain not found")
			}
			kestorePath := filepath.Join(a.homePath, "keystore", "wallets", k.chain)
			if err := os.MkdirAll(kestorePath, 0o755); err != nil {
				return err
			}
			addr, err := chain.ChainProvider.NewKeystore(k.password)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Created and encrypted: %s\n", addr)
			return nil
		},
	}
	k.chainFlag(new)
	k.passwordFlag(new, false)
	return new
}

// List keystore
func (k *keystoreState) list(a *appState) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := os.ReadDir(filepath.Join(a.homePath, "keystore", k.chain))
			if err != nil {
				return err
			}
			chain, err := a.config.Chains.Get(k.chain)
			if err != nil {
				return err
			}
			wallets := make(map[string]*types.Coin)
			for _, file := range files {
				name := file.Name()
				if strings.HasSuffix(name, ".pass") {
					wallet := strings.TrimSuffix(name, ".pass")
					balance, err := chain.ChainProvider.QueryBalance(cmd.Context(), wallet)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to query balance for %s: %v\n", wallet, err)
					}
					wallets[wallet] = balance
				}
			}
			printLabels("Wallet", "Balance")
			for wallet, balance := range wallets {
				if wallet == chain.ChainProvider.Config().GetWallet() {
					wallet = "* -> " + wallet
				}
				printValues(wallet, balance.Calculate())
			}
			return nil
		},
	}
	k.chainFlag(list)
	return list
}

// import keystore
func (k *keystoreState) importKey(a *appState) *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := a.config.Chains.Get(k.chain)
			if err != nil {
				return err
			}
			kestorePath := filepath.Join(a.homePath, "keystore", k.chain)
			if err := os.MkdirAll(kestorePath, 0o755); err != nil {
				return err
			}
			if _, err := os.Stat(k.path); os.IsNotExist(err) {
				return fmt.Errorf("file not found")
			}
			addr, err := chain.ChainProvider.ImportKeystore(cmd.Context(), k.path, k.password)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Key imported and Encrypted: %s\n", addr)
			return nil
		},
	}
	k.chainFlag(importCmd)
	k.passwordFlag(importCmd, false)
	k.keystorePathFlag(importCmd)
	return importCmd
}

// Use keystore using address
func (k *keystoreState) use(a *appState) *cobra.Command {
	use := &cobra.Command{
		Use:   "use",
		Short: "use keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := a.config.Chains.Get(k.chain)
			if err != nil {
				return err
			}
			kestorePath := filepath.Join(a.homePath, "keystore", k.chain, k.address)
			if _, err := os.Stat(kestorePath); os.IsNotExist(err) {
				return fmt.Errorf("keystore not found")
			}
			if _, err := os.Stat(kestorePath + ".pass"); os.IsNotExist(err) {
				return fmt.Errorf("password not found")
			}
			cf := chain.ChainProvider.Config()
			// check if it is the same wallet
			if cf.GetWallet() == k.address {
				fmt.Fprintf(os.Stdout, "Wallet already configured: %s\n", k.address)
				return nil
			}
			if err := chain.ChainProvider.SetAdmin(cmd.Context(), k.address); err != nil {
				return err
			}
			cf.SetWallet(k.address)
			if err := a.config.Save(a.configPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Wallet configured: %s\n", k.address)
			return nil
		},
	}
	k.chainFlag(use)
	k.addressFlag(use)
	return use
}

// change password
func (k *keystoreState) changePassword(a *appState) *cobra.Command {
	changePassword := &cobra.Command{
		Use:   "change-password",
		Short: "change password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if k.password != k.confirmPassword {
				return fmt.Errorf("password and confirm password does not match")
			}
			return fmt.Errorf("not implemented")
		},
	}
	k.chainFlag(changePassword)
	k.addressFlag(changePassword)
	k.passwordFlag(changePassword, true)
	return changePassword
}

// chain flag
func (k *keystoreState) chainFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&k.chain, "chain", "c", "", "chain id")
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
}

// password flag
func (k *keystoreState) passwordFlag(cmd *cobra.Command, isConfirmRequired bool) {
	cmd.Flags().StringVarP(&k.password, "password", "p", "", "password for keystore")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}
	if isConfirmRequired {
		cmd.Flags().StringVarP(&k.confirmPassword, "confirm", "c", "", "confirm password for keystore")
		if err := cmd.MarkFlagRequired("confirm"); err != nil {
			panic(err)
		}
	}
}

// address flag
func (k *keystoreState) addressFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&k.address, "address", "a", "", "address")
	if err := cmd.MarkFlagRequired("address"); err != nil {
		panic(err)
	}
}

func (k *keystoreState) keystorePathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&k.path, "keystore", "k", "", "keystore path")
	if err := cmd.MarkFlagRequired("keystore"); err != nil {
		panic(err)
	}
}
