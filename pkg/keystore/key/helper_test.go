package key_test

import (
	"context"
	"testing"

	pkg "github.com/scraly/go.common/pkg/keystore/key"

	"github.com/stretchr/testify/require"
)

func KeyGenerationGeneratorTest(generator pkg.Generator) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		key, err := generator(ctx)
		require.NoError(t, err, "Error should not be raised ongenration")
		require.NotNil(t, key, "Key holder should not be nil")
		require.NotEmpty(t, key.ID(), "Key identifier should not be empty")
		require.True(t, key.HasPrivate(), "Key should have a private key")
		require.True(t, key.HasPublic(), "Key should have a public key")

		pub := key.Public()
		require.NotNil(t, pub, "Public Key holder should not be nil")
		require.NotEmpty(t, pub.ID(), "Public Key identifier should not be empty")
		require.Equal(t, key.ID(), pub.ID(), "Key identifiers should be equals")
		require.False(t, pub.HasPrivate(), "Public Key should not have a private key")
		require.True(t, pub.HasPublic(), "Key should have a public key")
	}
}

func SignAndVerifyTest(generator pkg.Generator) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := []byte("toto")

		key, err := generator(ctx)
		require.NoError(t, err, "Error should not be raised ongenration")
		require.NotNil(t, key, "Key holder should not be nil")
		require.NotEmpty(t, key.ID(), "Key identifier should not be empty")
		require.True(t, key.HasPrivate(), "Key should have a private key")
		require.True(t, key.HasPublic(), "Key should have a public key")

		sig, err := key.Sign(data)
		require.NoError(t, err, "Error should not be raised on signature")
		require.NotNil(t, sig, "Signature should not be empty")

		err = key.Verify(data, sig)
		require.NoError(t, err, "Error should not be raised on verification")
	}
}

func KeyOperationsErrorTest(generator pkg.Generator) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		key, err := generator(ctx)
		require.NoError(t, err, "Error should not be raised ongenration")
		require.NotNil(t, key, "Key holder should not be nil")
		require.NotEmpty(t, key.ID(), "Key identifier should not be empty")
		require.True(t, key.HasPrivate(), "Key should have a private key")
		require.True(t, key.HasPublic(), "Key should have a public key")

		pub := key.Public()
		require.NotNil(t, pub, "Public Key holder should not be nil")
		require.NotEmpty(t, pub.ID(), "Public Key identifier should not be empty")
		require.Equal(t, key.ID(), pub.ID(), "Key identifiers should be equals")
		require.False(t, pub.HasPrivate(), "Public Key should not have a private key")
		require.True(t, pub.HasPublic(), "Key should have a public key")

		// Sign with public key
		sig, err := pub.Sign([]byte(""))
		require.Error(t, err, "Error should be raised when trying to sign with public key")
		require.Equal(t, pkg.ErrInvalidOperationCouldSignWithoutPrivateKey, err, "Error should be as expected")
		require.Empty(t, sig, "Signature should be empty on error")
	}
}
