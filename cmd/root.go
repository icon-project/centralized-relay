package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const appName = "centralized-relay"

var (
	defaultHome   = filepath.Join(os.Getenv("HOME"), ".centralized-relay")
	defaultDBName = "data"
	defaultConfig = "config.yaml"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd(nil)
	rootCmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt) // Using signal.Notify, instead of signal.NotifyContext, in order to see details of signal.
	go func() {
		// Wait for interrupt signal.
		sig := <-sigCh

		// Cancel context on root command.
		// If the invoked command respects this quickly, the main goroutine will quit right away.
		cancel()

		// Short delay before printing the received signal message.
		// This should result in cleaner output from non-interactive commands that stop quickly.
		time.Sleep(250 * time.Millisecond)
		fmt.Fprintf(os.Stderr, "Received signal %v. Attempting clean shutdown. Send interrupt again to force hard shutdown.\n", sig)

		// Dump all goroutines on panic, not just the current one.
		debug.SetTraceback("all")

		// Block waiting for a second interrupt or a timeout.
		// The main goroutine ought to finish before either case is reached.
		// But if a case is reached, panic so that we get a non-zero exit and a dump of remaining goroutines.
		select {
		case <-time.After(time.Minute):
			panic(errors.New("rly did not shut down within one minute of interrupt"))
		case sig := <-sigCh:
			panic(fmt.Errorf("received signal %v; forcing quit", sig))
		}
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// NewRootCmd returns the root command for relayer.
// If log is nil, a new zap.Logger is set on the app state
// based on the command line flags regarding logging.
func NewRootCmd(log *zap.Logger) *cobra.Command {
	// Use a local app state instance scoped to the new root command,
	// so that tests don't concurrently access the state.
	a := &appState{
		viper: viper.New(),
		log:   log,
	}

	// RootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "This application makes data relay between two chains!",
		Long:  strings.TrimSpace(`Use this to relay xcall packet between chains`),
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// Inside persistent pre-run because this takes effect after flags are parsed.
		if log == nil {
			log, err := newRootLogger(a.viper.GetString("log-format"), a.viper.GetBool("debug"))
			if err != nil {
				return err
			}

			a.log = log
		}

		if a.db == nil {
			db, err := lvldb.NewLvlDB(a.dbPath, false)
			if err != nil {
				return fmt.Errorf("error while creating db %v", err)
			}
			a.db = db
		}

		// reads `homeDir/config/config.yaml` into `a.Config`
		if err := a.loadConfigFile(rootCmd.Context()); err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, _ []string) {
		// Force syncing the logs before exit, if anything is buffered.
		_ = a.log.Sync()

		if a.db != nil {
			a.db.Close()
		}
	}

	// Register --home flag
	rootCmd.PersistentFlags().StringVar(&a.homePath, flagHome, defaultHome, "set home directory")
	if err := a.viper.BindPFlag(flagHome, rootCmd.PersistentFlags().Lookup(flagHome)); err != nil {
		panic(err)
	}

	// Register --debug flag
	rootCmd.PersistentFlags().BoolVarP(&a.debug, "debug", "d", false, "debug output")
	if err := a.viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().String("log-format", "auto", "log output format (auto, logfmt, json, or console)")
	if err := a.viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&a.configPath, "config-path", fmt.Sprintf("%s/%s", a.homePath, defaultConfig), "config path location")
	if err := a.viper.BindPFlag("config-path", rootCmd.PersistentFlags().Lookup("config-path")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&a.dbPath, "db-path", fmt.Sprintf("%s/%s", a.homePath, defaultDBName), "db path location")
	if err := a.viper.BindPFlag("db-path", rootCmd.PersistentFlags().Lookup("db-path")); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		startCmd(a),
		configCmd(a),
		chainsCmd(a),
		dbCmd(a),
		keystoreCmd(a),
	)
	return rootCmd
}

func newRootLogger(format string, debug bool) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006-01-02T15:04:05.000000Z07:00"))
	}
	config.LevelKey = "lvl"

	var enc zapcore.Encoder
	switch format {
	case "json":
		enc = zapcore.NewJSONEncoder(config)
	case "auto", "console":
		enc = zapcore.NewConsoleEncoder(config)
	case "logfmt":
		enc = zaplogfmt.NewEncoder(config)
	default:
		return nil, fmt.Errorf("unrecognized log format %q", format)
	}

	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	core := zapcore.NewTee(zapcore.NewCore(enc, os.Stderr, level))

	return zap.New(core), nil
}

// withUsage wraps a PositionalArgs to display usage only when the PositionalArgs
// variant is violated.
func withUsage(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := inner(cmd, args); err != nil {
			cmd.Root().SilenceUsage = false
			cmd.SilenceUsage = false
			return err
		}

		return nil
	}
}
