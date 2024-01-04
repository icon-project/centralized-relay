package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

type keystoreState struct {
	chain    string
	password string
	address  string
	client   *kms.Client
	app      *appState
}

func newKeyStoreState(ctx context.Context, app *appState) (*keystoreState, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile("default"))
	if err != nil {
		return nil, err
	}
	return &keystoreState{client: kms.NewFromConfig(cfg), app: app}, nil
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

	ks.AddCommand(state.init(), state.new())

	return ks
}

func (k *keystoreState) init() *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			input := &kms.CreateKeyInput{
				Description: aws.String("centralized-relay"),
			}
			output, err := k.client.CreateKey(cmd.Context(), input)
			if err != nil {
				return err
			}
			fmt.Println(output)
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
			kestorePath := fmt.Sprintf("%s/keystore/%s", k.app.homePath, k.chain)
			data, err := chain.ChainProvider.NewKeyStore(cmd.Context(), kestorePath, k.password)
			if err != nil {
				return err
			}
			address, err := chain.ChainProvider.AddressFromKeyStore(kestorePath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(fmt.Sprintf("%s/%s.json", kestorePath, address), data, 0o644); err != nil {
				return err
			}
			input := &kms.EncryptInput{
				KeyId:     &k.app.config.Global.KMSKeyID,
				Plaintext: []byte(k.password),
			}
			output, err := k.client.Encrypt(cmd.Context(), input)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, output.CiphertextBlob)
			return nil
		},
	}
	k.chainFlag(new)
	k.passwordFlag(new)
	return new
}

// import keystore
func (k *keystoreState) importKey() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			input := &kms.DecryptInput{
				KeyId:          &k.app.config.Global.KMSKeyID,
				CiphertextBlob: []byte(k.password),
			}
			output, err := k.client.Decrypt(cmd.Context(), input)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		},
	}
	k.chainFlag(importCmd)
	k.passwordFlag(importCmd)
	return importCmd
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

// Subcommand for keystore
// init Keystore
// new --chain=0x3 --password=1234
// use: --chain=0x3 --address=0x1234
// list --chain=0x3
// delete --chain=0x3 --address=0x1234
// import --chain=0x3 --keystore-path --password=1234
