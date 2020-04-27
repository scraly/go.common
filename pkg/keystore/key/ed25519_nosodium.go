// +build !sodium

package key

import (
	"golang.org/x/crypto/ed25519"
)

var (
	ed25519CryptoSignVerifyDetachedFunc = ed25519.Verify
	ed25519CryptoSignDetachedFunc       = ed25519.Sign
)
