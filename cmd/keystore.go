package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/spf13/cobra"
)

var TempDir = os.TempDir()

type keystoreState struct {
	client   kms.KMS
	chain    string
	password string
	address  string
	path     string
	app      *appState
}

func newKeyStoreState(ctx context.Context, app *appState) (*keystoreState, error) {
	k, err := kms.NewKMSConfig(ctx, "iconosphere")
	if err != nil {
		return nil, err
	}
	return &keystoreState{client: k, app: app}, nil
}

func keystoreCmd(a *appState) *cobra.Command {
	ks := &cobra.Command{
		Use:     "keystore",
		Aliases: []string{"ks"},
		Short:   "keystore utilty",
		Args:    withUsage(cobra.MaximumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s keystore [command]`, appName)),
	}
	state, err := newKeyStoreState(ks.Context(), a)
	if err != nil {
		panic(err)
	}

	ks.AddCommand(state.init(), state.new(), state.list(), state.importKey(), state.use())

	return ks
}

func (k *keystoreState) init() *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyID, err := k.client.Init(cmd.Context())
			if err != nil {
				return err
			}
			k.app.config.Global.KMSKeyID = *keyID
			if err := k.app.config.Save(k.app.homePath); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "KMS Key Created")
			fmt.Fprintln(os.Stdout, k.app.config.Global.KMSKeyID)
			return nil
		},
	}
	return init
}

func (k *keystoreState) new() *cobra.Command {
	new := &cobra.Command{
		Use:   "new",
		Short: "new keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, ok := k.app.config.Chains.GetAll()[k.chain]
			if !ok {
				return fmt.Errorf("chain not found")
			}
			kestorePath := filepath.Join(k.app.homePath, "keystore", k.chain)
			if err := os.MkdirAll(kestorePath, 0o755); err != nil {
				return err
			}
			addr, err := chain.ChainProvider.NewKeyStore(kestorePath, k.password)
			if err != nil {
				return err
			}
			data, err := k.client.Encrypt(cmd.Context(), k.app.config.Global.KMSKeyID, []byte(k.password))
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(kestorePath, fmt.Sprintf("%s.password", addr)), data, 0o644); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "KMS Key Encrypted")
			return nil
		},
	}
	k.chainFlag(new)
	k.passwordFlag(new)
	return new
}

// List keystore
func (k *keystoreState) list() *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := os.ReadDir(filepath.Join(k.app.homePath, "keystore", k.chain))
			if err != nil {
				return err
			}
			for _, file := range files {
				name := file.Name()
				if strings.HasSuffix(name, ".json") {
					fmt.Fprintln(os.Stdout, strings.TrimSuffix(name, ".json"))
				}
			}
			return nil
		},
	}
	k.chainFlag(list)
	return list
}

// import keystore
func (k *keystoreState) importKey() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	k.chainFlag(importCmd)
	k.passwordFlag(importCmd)
	k.keystorePathFlag(importCmd)
	return importCmd
}

// Use keystore using address
func (k *keystoreState) use() *cobra.Command {
	use := &cobra.Command{
		Use:   "use",
		Short: "use keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := k.app.config.Chains.Get(k.chain)
			if err != nil {
				return err
			}
			kestorePath := filepath.Join(k.app.homePath, "keystore", k.chain, k.address)
			if _, err := os.Stat(kestorePath + ".json"); os.IsNotExist(err) {
				return fmt.Errorf("keystore not found")
			}
			if _, err := os.Stat(kestorePath + ".password"); os.IsNotExist(err) {
				return fmt.Errorf("password not found")
			}
			cf := chain.ChainProvider.ProviderConfig()
			cf.SetWallet(k.address)
			if err := k.app.config.Save(k.app.homePath); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "Wallet configured")
			return nil
		},
	}
	k.chainFlag(use)
	k.addressFlag(use)
	return use
}

// chain flag
func (k *keystoreState) chainFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&k.chain, "chain", "c", "", "chain id")
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
}

// password flag
func (k *keystoreState) passwordFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&k.password, "password", "p", "", "password for keystore")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
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
	cmd.Flags().StringVarP(&k.path, "path", "k", "", "keystore path")
	if err := cmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}
}
