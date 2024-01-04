package kms

type KMS interface {
	Init(string, string) error
	Descrypt(data []byte) ([]byte, error)
	Encrypt(data []byte) ([]byte, error)
	Sign(data []byte) ([]byte, error)
	Verify(data []byte, sig []byte) (bool, error)
	Validate() error
}

type KMSConfig struct{}

func NewKMSConfig() *KMSConfig {
	return &KMSConfig{}
}

// Init creates a kms key for decrptying and encrypting data
func (c *KMSConfig) Init(secret string) error {
	// store the key in the aws ksm
}
