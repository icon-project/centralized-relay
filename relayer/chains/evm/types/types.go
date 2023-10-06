package types

type Metadata struct {
	Priority uint8
	Data     map[string]interface{}
}

type Message struct {
	Source      string
	Destination string
	Payload     []interface{}
	Metadata    Metadata
}
