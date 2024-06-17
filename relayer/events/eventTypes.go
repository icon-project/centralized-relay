package events

const (
	EmitMessage     = "emitMessage"
	CallMessage     = "callMessage"
	RollbackMessage = "rollBackMessage"

	// Special event types
	RevertMessage   = "revertMessage"
	SetAdmin        = "setAdmin"
	GetFee          = "getFee"
	SetFee          = "setFee"
	ClaimFee        = "claimFee"
	ExecuteRollback = "executeRollback"
)
