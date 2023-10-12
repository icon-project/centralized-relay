package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func parseMessageFromEventlog(log *zap.Logger, eventlogs []types.EventLog, height uint64) []providerTypes.Message {
	msgs := make([]providerTypes.Message, 0)
	return msgs
}
