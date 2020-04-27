package contexthelper

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestContextKey_Default(t *testing.T) {
	gomega.RegisterTestingT(t)

	const prefix string = "cCtxKeys"
	const keyName string = "testDefault"
	const complete string = prefix + keyName

	ref := Key{
		prefix:  prefix,
		keyName: keyName,
	}

	key := DefaultContextKey(keyName)

	gomega.Expect(key).Should(gomega.Equal(ref))
	gomega.Expect(key.String()).Should(gomega.Equal(complete))
}

func TestContextKey_CustomPrefix(t *testing.T) {
	gomega.RegisterTestingT(t)

	const prefix string = "MyCustomPrefix"
	const keyName string = "testCustom"
	const complete string = prefix + keyName

	ref := Key{
		prefix:  prefix,
		keyName: keyName,
	}

	key := ContextKey(prefix, keyName)

	gomega.Expect(key).Should(gomega.Equal(ref))
	gomega.Expect(key.String()).Should(gomega.Equal(complete))
}
