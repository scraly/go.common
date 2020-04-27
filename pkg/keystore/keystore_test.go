package keystore

import (
	"context"
	"crypto"
	"crypto/elliptic"
	"fmt"
	"testing"

	"github.com/hokaccha/go-prettyjson"

	"github.com/stretchr/testify/require"

	"github.com/scraly/go.common/pkg/keystore/backends/inmemory"
	"github.com/scraly/go.common/pkg/keystore/key"
)

func TestInMemoryKeyStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	backend, _ := inmemory.New()
	ks, _ := New(backend, WithInterval(1))

	// Generates keys
	k1, _ := ks.Generate(ctx, key.Ed25519)
	out, _ := prettyjson.Marshal(k1)
	fmt.Printf("%s\n", out)

	k2, _ := ks.Generate(ctx, key.ECDSA(elliptic.P256(), crypto.SHA256))
	out, _ = prettyjson.Marshal(k2)
	fmt.Printf("%s\n", out)

	// Add to keystore
	err := ks.Add(ctx, k1, k2.Public())
	require.NoError(t, err, "Error should be raised")

	// Expose all
	keys, err := ks.All(ctx)
	out, _ = prettyjson.Marshal(keys)
	fmt.Printf("%s\n", out)

	require.NoError(t, err, "Error should be raised")
	require.Equal(t, 2, len(keys), "Key count should be 2")

	// Get all keys
	for _, k := range keys {
		keyObject, _ := ks.Get(ctx, k.ID())

		if keyObject == nil {
			t.Fatal("Keystore : invalid key retrieval")
		}
	}

	// Get all public keys
	publicKeys, _ := ks.OnlyPublicKeys(ctx)
	for _, k := range publicKeys {
		keyObject, _ := ks.Get(ctx, k.ID())

		if keyObject == nil {
			t.Fatal("Keystore : invalid key retrieval")
		}
		if keyObject.HasPrivate() {
			t.Fatal("Keystore : should contain only public key in result")
		}
	}

	// Remove a key
	ks.Remove(ctx, k2.ID())
}
