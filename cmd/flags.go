package cmd

import (
	"github.com/icon-project/centralized-relay/relayer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHome            = "home"
	flagURL             = "url"
	flagJSON            = "json"
	flagYAML            = "yaml"
	flagOverride        = "override"
	flagVersion         = "version"
	flagDebugAddr       = "debug-addr"
	flagOverwriteConfig = "overwrite"
	flagFlushInterval   = "flush-interval"
	flagFresh           = "fresh"
)

func flushIntervalFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().DurationP(flagFlushInterval, "i", relayer.DefaultFlushInterval, "how frequently should a flush routine be run")
	if err := v.BindPFlag(flagFlushInterval, cmd.Flags().Lookup(flagFlushInterval)); err != nil {
		panic(err)
	}
	return cmd
}
func freshFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Bool(flagFresh, false, "whether to clear db and tart fresh")
	if err := v.BindPFlag(flagFresh, cmd.Flags().Lookup(flagFresh)); err != nil {
		panic(err)
	}
	return cmd
}
