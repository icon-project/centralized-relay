package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func chainsCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chains",
		Aliases: []string{"ch"},
		Short:   "Manage chain configurations",
	}

	cmd.AddCommand(
		chainsListCmd(a),
		chainsAddCmd(a),
		chainsDeleteCmd(a),
	)

	return cmd
}

func chainsListCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Returns chain configuration data",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s chains list
$ %s ch l`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsn, err := cmd.Flags().GetBool(flagJSON)
			if err != nil {
				return err
			}

			yml, err := cmd.Flags().GetBool(flagYAML)
			if err != nil {
				return err
			}

			configs := a.config.Wrapped().ProviderConfigs
			if len(configs) == 0 {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: no chains found (do you need to run 'rly chains add'?)")
			}

			switch {
			case yml && jsn:
				return fmt.Errorf("can't pass both --json and --yaml, must pick one")
			case yml:
				out, err := yaml.Marshal(configs)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			case jsn:
				out, err := jsoniter.Marshal(configs)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			default:
				i := 0
				for _, c := range a.config.Chains {
					i++
					fmt.Fprintf(cmd.OutOrStdout(), "%d: %-20s -> type(%s)\n", i, c.NID(), c.ChainProvider.Type())
				}
				return nil
			}
		},
	}
	return yamlFlag(a.viper, jsonFlag(a.viper, cmd))
}

func chainsAddCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain-name...]",
		Aliases: []string{"a"},
		Short: "Add a new chain to the configuration file by fetching chain metadata from \n" +
			" passing a file (-f) ",
		Args: withUsage(cobra.MinimumNArgs(0)),
		Example: fmt.Sprintf(` $ %s chains add cosmoshub
 $ %s chains add --file chains/ibc0.json ibc0`, appName, appName),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := cmd.Flags().GetString(flagFile)
			if err != nil {
				return fmt.Errorf("File is not present")
			}

			if ok := a.config; ok == nil {
				return fmt.Errorf("config not initialized, consider running `rly config init`")
			}

			return a.performConfigLockingOperation(cmd.Context(), func() error {
				// default behavior fetch from chain registry
				// still allow for adding config from url or file
				switch {
				case file != "":
					var chainName string
					switch len(args) {
					case 0:
						chainName = strings.Split(filepath.Base(file), ".")[0]
					case 1:
						chainName = args[0]
					default:
						return errors.New("one chain name is required")
					}
					if err := addChainFromFile(a, chainName, file); err != nil {
						return err
					}

				default:
					return fmt.Errorf("file not present")
				}
				return nil
			})
		},
	}

	cmd = fileFlag(a.viper, cmd)
	return cmd
}

func addChainFromFile(a *appState, chainName string, file string) error {
	// If the user passes in a file, attempt to read the chain config from that file
	var pcw ProviderConfigWrapper
	if _, err := os.Stat(file); err != nil {
		return err
	}

	byt, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err = jsoniter.Unmarshal(byt, &pcw); err != nil {
		return err
	}

	prov, err := pcw.Value.NewProvider(context.Background(),
		a.log.With(zap.String("provider_type", pcw.Type)),
		a.homePath, a.debug, chainName,
	)
	if err != nil {
		return fmt.Errorf("failed to build ChainProvider for %s: %w", file, err)
	}

	c := relayer.NewChain(a.log, prov, a.debug)
	if err = a.config.AddChain(c); err != nil {
		return err
	}

	return nil
}

func chainsDeleteCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete chain_name",
		Aliases: []string{"d"},
		Short:   "Removes chain from config based off chain-id",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s chains delete ibc-0
$ %s ch d ibc-0`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain := args[0]
			return a.performConfigLockingOperation(cmd.Context(), func() error {
				_, ok := a.config.Chains[chain]
				if !ok {
					return errChainNotFound(chain)
				}
				a.config.DeleteChain(chain)
				return nil
			})
		},
	}
	return cmd
}

func (c *Config) DeleteChain(chain string) {
	delete(c.Chains, chain)
}
