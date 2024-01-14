package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func newKeyStoreState(ctx context.Context) (*keystoreState, error) {
	return &keystoreState{}, nil
}

func keystoreCmd(a *appState) *cobra.Command {
	ks := &cobra.Command{
		Use:     "keystore",
		Aliases: []string{"ks"},
		Short:   "keystore utilty",
		Args:    withUsage(cobra.MaximumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s keystore [command]`, appName)),
	}
	state, err := newKeyStoreState(ks.Context())
	if err != nil {
		panic(err)
	}

	ks.AddCommand(state.init(a), state.new(a), state.list(a), state.importKey(a), state.use(a))

	return ks
}

func (k *keystoreState) init(a *appState) *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.kms.Init(cmd.Context()); err != nil {
				return err
			}
			if err := a.config.Save(a.homePath); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "KMS Key Created")
			fmt.Fprintln(os.Stdout, a.config.Global.KMSKeyID)
			return nil
		},
	}
	return init
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
			kestorePath := filepath.Join(a.homePath, "keystore", k.chain)
			if err := os.MkdirAll(kestorePath, 0o755); err != nil {
				return err
			}
			addr, err := chain.ChainProvider.NewKeyStore(kestorePath, k.password)
			if err != nil {
				return err
			}
			data, err := a.kms.Encrypt(cmd.Context(), []byte(k.password))
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
				return fmt.Errorf("keystore not found")
			}
			addr, err := chain.ChainProvider.AddressFromKeyStore(k.path, k.password)
			if err != nil {
				return err
			}
			data, err := os.Open(k.path)
			if err != nil {
				return err
			}
			defer data.Close()
			keystore, err := os.OpenFile(filepath.Join(kestorePath, fmt.Sprintf("%s.json", addr)), os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return err
			}
			defer keystore.Close()
			if _, err := io.Copy(keystore, data); err != nil {
				return err
			}
			secret, err := a.kms.Encrypt(cmd.Context(), []byte(k.password))
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(kestorePath, fmt.Sprintf("%s.password", addr)), secret, 0o644); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "KMS Key imported and Encrypted")
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
			if _, err := os.Stat(kestorePath + ".json"); os.IsNotExist(err) {
				return fmt.Errorf("keystore not found")
			}
			if _, err := os.Stat(kestorePath + ".password"); os.IsNotExist(err) {
				return fmt.Errorf("password not found")
			}
			cf := chain.ChainProvider.ProviderConfig()
			cf.SetWallet(k.address)
			if err := a.config.Save(a.homePath); err != nil {
				return err
			}
			// TODO: set admin
			// if err := chain.ChainProvider.SetAdmin(cmd.Context(), k.address); err != nil {
			// 	return err
			// }
			fmt.Fprintln(os.Stdout, "Wallet configured")
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
			// TODO: implement change password
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
		cmd.Flags().StringVarP(&k.password, "confirm", "c", "", "confirm password for keystore")
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
