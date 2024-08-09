package events

const (
	EmitMessage     = "emitMessage"
	CallMessage     = "callMessage"
	RollbackMessage = "rollbackMessage"

	// Special event types
	RevertMessage = "revertMessage"
	SetAdmin      = "setAdmin"
	GetFee        = "getFee"
	SetFee        = "setFee"
	ClaimFee      = "claimFee"
)
