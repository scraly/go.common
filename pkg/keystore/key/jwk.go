package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ed25519"
)

// rawJWK implements the internal representation for serialzing/deserializing a JWK: RFC 7517 Section 4
type rawJWK struct {
	IssuedAt                 int64    `json:"iat,omitempty"`
	PublicKeyUse             string   `json:"use,omitempty"`      // JWK 4.2
	KeyType                  string   `json:"kty,omitempty"`      // JWK 4.1
	KeyID                    string   `json:"kid,omitempty"`      // JWK 4.5
	KeyOperations            []string `json:"key_ops,omitempty"`  // JWK 4.3
	Curve                    string   `json:"crv,omitempty"`      // RSA Curve JWA 6.2.1.1
	Algorithm                string   `json:"alg,omitempty"`      // JWK 4.4
	K                        string   `json:"k,omitempty"`        // Symmetric Key JWA 6.4.1
	X                        string   `json:"x,omitempty"`        // RSA X Coordindate JWA 6.2.1.2
	Y                        string   `json:"y,omitempty"`        // RSA Y Coordinate JWA 6.2.1.3
	N                        string   `json:"n,omitempty"`        // RSA Modulus, JWA 6.3.1.1
	E                        string   `json:"e,omitempty"`        // RSA Exponent JWA 6.3.1.2
	D                        string   `json:"d,omitempty"`        // RSA Private Exponent JWA 6.3.2.1, ECC Private Key JWA 6.2.2.1
	P                        string   `json:"p,omitempty"`        // RSA First Prime Factor JWA 6.3.2.2
	Q                        string   `json:"q,omitempty"`        // RSA Second Prime Factor JWA 6.3.2.3
	Dp                       string   `json:"dp,omitempty"`       // RSA First Factor CRT Exponent JWA 6.3.2.4
	Dq                       string   `json:"dq,omitempty"`       // RSA SEcond Factor CRT Exponent JWA 6.3.2.5
	Qi                       string   `json:"qi,omitempty"`       // RSA First CRT Coefficient JWA 6.3.2.6
	X509URL                  string   `json:"x5u,omitempty"`      // JWK 4.6
	X509CertChain            []string `json:"x5c,omitempty"`      // JWK 4.7
	X509Sha1Thumbprint       string   `json:"x5t,omitempty"`      // JWK 4.8
	X509CertSha256Thumbprint string   `json:"x5t#S256,omitempty"` // JWK 4.9
}

// -----------------------------------------------------------------------------

func toEd25519(raw *rawJWK) (Key, error) {
	x, err := base64.RawURLEncoding.DecodeString(raw.X)
	if err != nil {
		return nil, err
	}

	if len(x) != ed25519.PublicKeySize {
		return nil, errors.New("key: invalid ed25519 public key size")
	}

	k := &ed25519Key{
		pub: x,
	}
	if len(raw.D) > 0 {
		d, err := base64.RawURLEncoding.DecodeString(raw.D)
		if err != nil {
			return nil, err
		}
		if len(d) != ed25519.PrivateKeySize {
			return nil, errors.New("key: invalid ed25519 private key size")
		}
		k.priv = d
	}

	return k, nil
}

func toECDSA(raw *rawJWK) (Key, error) {
	if raw.Curve == "" || raw.X == "" || raw.Y == "" {
		return nil, errors.New("key: malformed JWK EC key")
	}

	var curve elliptic.Curve
	switch raw.Curve {
	case "P-224":
		curve = elliptic.P224()
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("Unknown curve type: %s", raw.Curve)
	}

	pubKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     &big.Int{},
		Y:     &big.Int{},
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(raw.X)
	if err != nil {
		return nil, fmt.Errorf("key: malformed JWK EC key, %s", err)
	}
	pubKey.X.SetBytes(xBytes)

	yBytes, err := base64.RawURLEncoding.DecodeString(raw.Y)
	if err != nil {
		return nil, fmt.Errorf("key: malformed JWK EC key, %s", err)
	}
	pubKey.Y.SetBytes(yBytes)

	key := &ecdsaKey{
		curve: curve,
		kid:   raw.KeyID,
		pub:   pubKey,
	}

	if len(raw.D) > 0 {
		privKey := &ecdsa.PrivateKey{
			PublicKey: *pubKey,
			D:         &big.Int{},
		}

		dBytes, err := base64.RawURLEncoding.DecodeString(raw.D)
		if err != nil {
			return nil, fmt.Errorf("key: malformed JWK EC key, %s", err)
		}
		privKey.D.SetBytes(dBytes)

		key.priv = privKey
	}

	return key, nil
}
