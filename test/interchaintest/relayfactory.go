package interchaintest

import (
	"testing"

	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"go.uber.org/zap"
)

const (
	rlyRelayerUser = "100:1000"
)

type RelayerFactory interface {
	// Build returns a Relayer associated with the given arguments.
	Build(t *testing.T, cli *client.Client, networkID string) ibc.Relayer

	// Name returns a descriptive name of the factory,
	// indicating details of the Relayer that will be built.
	Name() string

	// Capabilities is an indication of the features this relayer supports.
	// Tests for any unsupported features will be skipped rather than failed.
	Capabilities() map[relayer.Capability]bool
}

// builtinRelayerFactory is the built-in relayer factory that understands
// how to start the cosmos relayer in a docker container.
type builtinRelayerFactory struct {
	log     *zap.Logger
	options relayer.RelayerOptions
}

func NewICONRelayerFactory(logger *zap.Logger, options ...relayer.RelayerOption) RelayerFactory {
	return builtinRelayerFactory{log: logger, options: options}
}

// Build returns a relayer chosen depending on f.impl.
func (f builtinRelayerFactory) Build(t *testing.T, cli *client.Client, networkID string) ibc.Relayer {
	return centralized.NewCentralizedRelayer(f.log, t.Name(), cli, networkID, f.options...)
}

func (f builtinRelayerFactory) Name() string {
	return "iconRelayer@"
}

// Capabilities returns the set of capabilities for the
// relayer implementation backing this factory.
func (f builtinRelayerFactory) Capabilities() map[relayer.Capability]bool {
	return centralized.Capabilities()
}

// Config holds configuration values for the relayer used in the tests.
type Config struct {
	// Tag is the tag used for the relayer image.
	Tag string `mapstructure:"tag"`
	// Image is the image that should be used for the relayer.
	Image string `mapstructure:"image"`
	// KMS_ID is the kms is that should be used for the relayer.
	KMS_ID string `mapstructure:"kms_id"`
	// KMS_URL is the kms endpoint that should be used for the relayer.
	KMS_URL string `mapstructure:"kms_url"`
}

// New returns an implementation of ibc.Relayer depending on the provided RelayerType.
func New(t *testing.T, cfg Config, logger *zap.Logger, dockerClient *client.Client, network string) ibc.Relayer {
	optionDocker := relayer.CustomDockerImage(cfg.Image, cfg.Tag, rlyRelayerUser)
	//flagOptions := relayer.StartupFlags("-p", "events") // relayer processes via events
	imageOptions := relayer.ImagePull(false)

	relayerFactory := NewICONRelayerFactory(logger, optionDocker, imageOptions)
	return relayerFactory.Build(t, dockerClient, network)
}
