// +build sodium

// #cgo pkg-config: libsodium

package key

import "github.com/GoKillers/libsodium-go/cryptosign"

func ed25519CryptoSignVerifyDetachedFunc(publicKey PublicKey, message, sig []byte) bool {
	ok := cryptosign.CryptoSignVerifyDetached(sig, message, publicKey)
	return ok > 0
}

func ed25519CryptoSignDetachedFunc(privateKey PrivateKey, message []byte) []byte {
	sig, ok := cryptosign.CryptoSignDetached(message, privateKey)
	if ok < 0 {
		panic("Unable to sign message with libsodium")
	}
	return sig
}
