package icon

var (
	// Connection Methods
	MethodSendMessage   = "sendMessage"
	MethodRecvMessage   = "recvMessage"
	MethodGetReceipts   = "getReceipts"
	MethodSetAdmin      = "setAdmin"
	MethodRevertMessage = "revertMessage"
	MethodGetFee        = "getFee"
	MethodSetFee        = "setFee"
	MethodClaimFees     = "claimFees"

	// XCALL Methods
	MethodExecuteCall     = "executeCall"
	MethodExecuteRollback = "executeRollback"

	// Cluster Methods
	MethodSubmitPacket             = "submitPacket"
	MethodClusterMsgReceived       = "packetSubmitted"
	MethodPacketAcknowledged       = "packetAcknowledged"
	MethodRecvMessageWithSignature = "recvMessageWithSignatures"
)
