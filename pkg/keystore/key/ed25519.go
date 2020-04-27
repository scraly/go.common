package key

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/dchest/uniuri"
	"golang.org/x/crypto/ed25519"
)

type ed25519Key struct {
	timestamp time.Time
	kid       string
	priv      ed25519.PrivateKey
	pub       ed25519.PublicKey
}

// Ed25519 key holder
func Ed25519(context.Context) (Key, error) {

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ed25519Key{
		kid:       uniuri.NewLen(12),
		timestamp: time.Now().UTC(),
		priv:      privateKey,
		pub:       publicKey,
	}, nil
}

// -----------------------------------------------------------------------------

func (k *ed25519Key) ID() string {
	return k.kid
}

func (k *ed25519Key) HasPrivate() bool {
	return len(k.priv) > 0
}

func (k *ed25519Key) HasPublic() bool {
	return len(k.pub) > 0
}

func (k *ed25519Key) Public() Key {
	return &ed25519Key{
		kid: k.ID(),
		pub: k.pub,
	}
}

func (k *ed25519Key) Sign(data []byte) ([]byte, error) {
	if !k.HasPrivate() {
		return nil, ErrInvalidOperationCouldSignWithoutPrivateKey
	}

	sig := ed25519CryptoSignDetachedFunc(k.priv, data)
	return sig, nil
}

func (k *ed25519Key) Verify(data, sig []byte) error {
	if !k.HasPublic() {
		return ErrInvalidOperationCouldVerifyWithoutPublicKey
	}

	if valid := ed25519CryptoSignVerifyDetachedFunc(k.pub, data, sig); !valid {
		return ErrInvalidSignature
	}

	return nil
}

// -----------------------------------------------------------------------------

func (k *ed25519Key) MarshalJSON() ([]byte, error) {
	r := &rawJWK{
		KeyID:         k.ID(),
		KeyType:       "OKP",
		Algorithm:     "EC",
		Curve:         "Ed25519",
		X:             base64.RawURLEncoding.EncodeToString(k.pub),
		PublicKeyUse:  "sig",
		KeyOperations: []string{"verify"},
	}
	if k.HasPrivate() {
		r.D = base64.RawURLEncoding.EncodeToString(k.priv)
		r.KeyOperations = append(r.KeyOperations, "sign")
	}

	return json.Marshal(r)
}
