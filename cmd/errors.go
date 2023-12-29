package cmd

import "fmt"

func errChainNotFound(chainName string) error {
	return fmt.Errorf("chain with name \"%s\" not found in config. consider running `rly chains add %s`", chainName, chainName)
}
