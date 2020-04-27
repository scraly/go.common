package key

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

type ecdsaKey struct {
	timestamp time.Time
	kid       string
	priv      *ecdsa.PrivateKey
	pub       *ecdsa.PublicKey

	alg   crypto.Hash
	curve elliptic.Curve
}

// ECDSA key holder
func ECDSA(curve elliptic.Curve, alg crypto.Hash) func(context.Context) (Key, error) {
	return func(ctx context.Context) (Key, error) {

		// Generate a RSA keypair
		privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, err
		}

		// Extract public key
		publicKey := privateKey.Public().(*ecdsa.PublicKey)

		// Create the hasher
		if !alg.Available() {
			return nil, ErrAlgorithmNotSupported
		}

		return &ecdsaKey{
			kid:       uniuri.NewLen(12),
			timestamp: time.Now().UTC(),
			priv:      privateKey,
			pub:       publicKey,
			alg:       alg,
			curve:     curve,
		}, nil
	}
}

// -----------------------------------------------------------------------------

func (k *ecdsaKey) ID() string {
	return k.kid
}

func (k *ecdsaKey) HasPrivate() bool {
	return k.priv != nil
}

func (k *ecdsaKey) HasPublic() bool {
	return k.pub != nil
}

func (k *ecdsaKey) Public() Key {
	return &ecdsaKey{
		kid: k.ID(),
		pub: k.pub,
	}
}

// -----------------------------------------------------------------------------

func (k *ecdsaKey) Sign(data []byte) ([]byte, error) {
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

	keysiz := k.alg.Size()
	curveBits := k.curve.Params().BitSize
	if curveBits != keysiz*8 {
		return nil, errors.New("key size does not match curve bit size")
	}

	// Sign the string and return r, s
	r, s, err := ecdsa.Sign(rand.Reader, k.priv, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}
	out := make([]byte, keysiz*2)
	copy(out, r.Bytes())
	copy(out[keysiz:], s.Bytes())

	return out, nil
}

func (k *ecdsaKey) Verify(data, sig []byte) error {
	if !k.HasPublic() {
		return ErrInvalidOperationCouldVerifyWithoutPublicKey
	}

	// Hash given data
	hasher := k.alg.New()
	if _, err := hasher.Write(data); err != nil {
		return errors.Wrap(err, "Unable to hash given data")
	}

	curveOrderByteSize := k.pub.Curve.Params().P.BitLen() / 8
	if len(sig) < curveOrderByteSize {
		return ErrInvalidSignature
	}

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sig[:curveOrderByteSize])
	s.SetBytes(sig[curveOrderByteSize:])

	// Verify the signature
	if verifystatus := ecdsa.Verify(k.pub, hasher.Sum(nil), r, s); !verifystatus {
		return ErrInvalidSignature
	}

	// No error
	return nil
}

// -----------------------------------------------------------------------------

func (k *ecdsaKey) MarshalJSON() ([]byte, error) {
	r := &rawJWK{
		KeyID:         k.ID(),
		KeyType:       "EC",
		Curve:         k.pub.Params().Name,
		X:             base64.RawURLEncoding.EncodeToString(k.pub.X.Bytes()),
		Y:             base64.RawURLEncoding.EncodeToString(k.pub.Y.Bytes()),
		PublicKeyUse:  "sig",
		KeyOperations: []string{"verify"},
	}
	if k.HasPrivate() {
		r.D = base64.RawURLEncoding.EncodeToString(k.priv.D.Bytes())
		r.KeyOperations = append(r.KeyOperations, "sign")
	}

	return json.Marshal(r)
}
