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
	flagFile            = "file"
	flagConfig          = "config"
)

func flushIntervalFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().DurationP(flagFlushInterval, "i", relayer.DefaultFlushInterval, "how frequently should a flush routine be run")
	if err := v.BindPFlag(flagFlushInterval, cmd.Flags().Lookup(flagFlushInterval)); err != nil {
		panic(err)
	}
	return cmd
}

func freshFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Bool(flagFresh, false, "whether to clear db and start fresh")
	if err := v.BindPFlag(flagFresh, cmd.Flags().Lookup(flagFresh)); err != nil {
		panic(err)
	}
	return cmd
}

func yamlFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagYAML, "y", false, "output using yaml")
	if err := v.BindPFlag(flagYAML, cmd.Flags().Lookup(flagYAML)); err != nil {
		panic(err)
	}
	return cmd
}

func jsonFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagJSON, "j", false, "returns the response in json format")
	if err := v.BindPFlag(flagJSON, cmd.Flags().Lookup(flagJSON)); err != nil {
		panic(err)
	}
	return cmd
}

func fileFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagFile, "f", "", "fetch json data from specified file")
	if err := v.BindPFlag(flagFile, cmd.Flags().Lookup(flagFile)); err != nil {
		panic(err)
	}
	return cmd
}
