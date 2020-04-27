package key

import "errors"

// Key contract for key information holder
type Key interface {
	ID() string
	HasPrivate() bool
	HasPublic() bool
	Public() Key
	Sign(data []byte) ([]byte, error)
	Verify(data []byte, sig []byte) error
}

// -----------------------------------------------------------------------------

var (
	// ErrInvalidSignature is raised when the signature could not be verified
	ErrInvalidSignature = errors.New("key: invalid signature")
	// ErrInvalidOperationCouldSignWithoutPrivateKey is raised when trying to sign using the public key
	ErrInvalidOperationCouldSignWithoutPrivateKey = errors.New("key: invalid operation : could not sign without a private key")
	// ErrInvalidOperationCouldVerifyWithoutPublicKey is raised when trying to verify signature without matching public key
	ErrInvalidOperationCouldVerifyWithoutPublicKey = errors.New("key: invalid operation : could not verify without a public key")
	// ErrAlgorithmNotSupported is raised when using not load algorithm (missing imports)
	ErrAlgorithmNotSupported = errors.New("key: Algorithm not supported")
)
