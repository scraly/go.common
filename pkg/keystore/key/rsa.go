package key

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

type rsaKey struct {
	timestamp time.Time
	kid       string
	priv      *rsa.PrivateKey
	pub       *rsa.PublicKey

	alg crypto.Hash
}

// RSA key holder
func RSA(size int, alg crypto.Hash) func(context.Context) (Key, error) {
	return func(ctx context.Context) (Key, error) {

		// Generate a RSA keypair
		privateKey, err := rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			return nil, err
		}
		privateKey.Precompute()

		// Extract public key
		publicKey := privateKey.Public().(*rsa.PublicKey)

		return &rsaKey{
			kid:       uniuri.NewLen(12),
			timestamp: time.Now().UTC(),
			priv:      privateKey,
			pub:       publicKey,
			alg:       alg,
		}, nil
	}
}

// -----------------------------------------------------------------------------

func (k *rsaKey) ID() string {
	return k.kid
}

func (k *rsaKey) HasPrivate() bool {
	return k.priv != nil
}

func (k *rsaKey) HasPublic() bool {
	return k.pub != nil
}

func (k *rsaKey) Public() Key {
	return &rsaKey{
		kid: k.ID(),
		pub: k.pub,
	}
}

// -----------------------------------------------------------------------------

func (k *rsaKey) Sign(data []byte) ([]byte, error) {
	if !k.HasPrivate() {
		return nil, ErrInvalidOperationCouldSignWithoutPrivateKey
	}

	// Create the hasher
	if !k.alg.Available() {
		return nil, ErrAlgorithmNotSupported
	}
	hasher := k.alg.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to hash given data")
	}

	// Return encoded signature
	return rsa.SignPKCS1v15(rand.Reader, k.priv, k.alg, hasher.Sum(nil))
}

func (k *rsaKey) Verify(data, sig []byte) error {
	if !k.HasPublic() {
		return ErrInvalidOperationCouldVerifyWithoutPublicKey
	}

	// Create the hasher
	if !k.alg.Available() {
		return ErrAlgorithmNotSupported
	}
	hasher := k.alg.New()
	_, err := hasher.Write(data)
	if err != nil {
		return errors.Wrap(err, "Unable to hash given data")
	}

	// Verify the signature
	return rsa.VerifyPKCS1v15(k.pub, k.alg, hasher.Sum(nil), sig)
}

// -----------------------------------------------------------------------------

func (k *rsaKey) MarshalJSON() ([]byte, error) {
	r := &rawJWK{
		KeyID:         k.ID(),
		KeyType:       "RSA",
		N:             base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.pub.E)).Bytes()),
		E:             base64.RawURLEncoding.EncodeToString(k.pub.N.Bytes()),
		PublicKeyUse:  "sig",
		KeyOperations: []string{"verify"},
	}
	if k.HasPrivate() {
		r.D = base64.RawURLEncoding.EncodeToString(k.priv.D.Bytes())
		r.KeyOperations = append(r.KeyOperations, "sign")
	}

	return json.Marshal(r)
}
