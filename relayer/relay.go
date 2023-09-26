package relayer

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	DefaultFlushInterval = 5 * time.Minute
)

// main start loop
func Start(
	ctx context.Context,
	log *zap.Logger,
	chains map[string]*Chain,
	flushInterval time.Duration,
	fresh bool,
) chan error {

	errorChan := make(chan error, 1)

	fmt.Println("started main loop")

	go Relay(
		ctx,
		log,
		chains,
		flushInterval,
		fresh,
		errorChan,
	)

	return errorChan
}

func Relay(
	ctx context.Context,
	log *zap.Logger,
	chains map[string]*Chain,
	flushInterval time.Duration,
	fresh bool,
	errCh chan error,
) {

}
