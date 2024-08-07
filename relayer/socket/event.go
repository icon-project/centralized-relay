package socket

const (
	EventGetBlock       Event = "GetBlock"
	EventGetMessageList Event = "GetMessageList"
	EventRelayMessage   Event = "RelayMessage"
	EventMessageRemove  Event = "MessageRemove"
	EventPruneDB        Event = "PruneDB"
	EventRevertMessage  Event = "RevertMessage"
	EventError          Event = "Error"
	EventGetFee         Event = "GetFee"
	EventSetFee         Event = "SetFee"
	EventClaimFee       Event = "ClaimFee"
	EventCurrentHeight  Event = "CurrentHeight"
	EventChainConfig    Event = "ChainConfig"
)
