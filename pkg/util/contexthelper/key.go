package contexthelper

// DefaultContextKey ...
func DefaultContextKey(keyName string) Key {
	return ContextKey("cCtxKeys", keyName)
}

// ContextKey ...
func ContextKey(prefix string, keyName string) Key {
	return Key{
		prefix:  prefix,
		keyName: keyName,
	}
}

// Key is used to generate keys for Context value
type Key struct {
	prefix  string
	keyName string
}

func (c Key) String() string {
	return c.prefix + c.keyName
}
