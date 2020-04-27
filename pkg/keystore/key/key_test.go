package key_test

import (
	"context"
	"crypto"
	"crypto/elliptic"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	_ "crypto/sha256"
	_ "crypto/sha512"

	pkg "github.com/scraly/go.common/pkg/keystore/key"
)

var (
	keyGenerators = map[string]pkg.Generator{
		"ed25519":        pkg.Ed25519,
		"ecp256-sha256":  pkg.ECDSA(elliptic.P256(), crypto.SHA256),
		"rsa2048-sha256": pkg.RSA(2048, crypto.SHA256),
	}
)

func TestKeyGeneration(t *testing.T) {
	for k, generator := range keyGenerators {
		t.Run(fmt.Sprintf("case=%s", k), KeyGenerationGeneratorTest(generator))
	}
}

func TestSignAndVerify(t *testing.T) {
	for k, generator := range keyGenerators {
		t.Run(fmt.Sprintf("case=%s", k), SignAndVerifyTest(generator))
	}
}

func TestKeyOperationError(t *testing.T) {
	for k, generator := range keyGenerators {
		t.Run(fmt.Sprintf("case=%s", k), KeyOperationsErrorTest(generator))
	}
}

func BenchmarkED25519KeySignature(b *testing.B) {
	key, err := pkg.Ed25519(context.Background())
	require.NoError(b, err)

	data := []byte("toto")
	for n := 0; n < b.N; n++ {
		_, err := key.Sign(data)
		require.NoError(b, err)
	}
}

func BenchmarkEcdsaKeySignature(b *testing.B) {
	key, err := pkg.ECDSA(elliptic.P256(), crypto.SHA256)(context.Background())
	require.NoError(b, err)

	data := []byte("toto")
	for n := 0; n < b.N; n++ {
		_, err := key.Sign(data)
		require.NoError(b, err)
	}
}

func BenchmarkRSAKeySignature(b *testing.B) {
	key, err := pkg.RSA(2048, crypto.SHA256)(context.Background())
	require.NoError(b, err)

	data := []byte("toto")
	for n := 0; n < b.N; n++ {
		_, err := key.Sign(data)
		require.NoError(b, err)
	}
}
