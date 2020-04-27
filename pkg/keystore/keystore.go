package keystore

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/scraly/go.common/pkg/keystore/backends"
	"github.com/scraly/go.common/pkg/keystore/key"
	"github.com/scraly/go.common/pkg/log"
	"go.uber.org/zap"

	"github.com/golang/snappy"
	"github.com/pkg/errors"
)

type defaultKeyStore struct {
	sync.RWMutex

	store backends.Backend
	dopts *Options
	done  context.CancelFunc

	keys map[string]key.Key
}

// New returns a default keystore implementation instance
func New(backend backends.Backend, opts ...Option) (KeyStore, error) {
	// Default Options
	options := &Options{
		Interval: 60,
		Watch:    false,
		OneTime:  false,
		Snappy:   true,
	}

	// Overrides with option
	for _, opt := range opts {
		opt(options)
	}

	// Initializes default value
	ks := &defaultKeyStore{
		store: backend,
		dopts: options,
		keys:  make(map[string]key.Key),
	}

	// Synchronize with backend
	return ks, nil
}

// -----------------------------------------------------------------------------
func (ks *defaultKeyStore) Generate(ctx context.Context, generator key.Generator) (key.Key, error) {
	k, err := generator(ctx)
	if err != nil {
		return nil, fmt.Errorf("keystore: Key generation error %v", err)
	}

	return k, nil
}

func (ks *defaultKeyStore) All(_ context.Context) ([]key.Key, error) {
	ks.RLock()
	defer ks.RUnlock()

	var result []key.Key
	for _, i := range ks.keys {
		result = append(result, i)
	}

	return result, nil
}

func (ks *defaultKeyStore) OnlyPublicKeys(_ context.Context) ([]key.Key, error) {
	ks.RLock()
	defer ks.RUnlock()

	var result []key.Key
	for _, i := range ks.keys {
		result = append(result, i.Public())
	}

	return result, nil
}

func (ks *defaultKeyStore) Add(ctx context.Context, keys ...key.Key) error {
	for _, k := range keys {
		// Encode key
		jwk, err := json.Marshal(k.Public())
		if err != nil {
			continue
		}

		if ks.dopts.Snappy {
			// Compress JWK
			jwk = snappy.Encode(nil, jwk)
		}

		// Wrap the key in a holder
		holder := &keyHolder{
			IssuedAt: time.Now().UTC().Unix(),
			Data:     jwk,
		}

		// Marshal only public key to json
		payload, err := json.Marshal(holder)
		if err != nil {
			return fmt.Errorf("keystore: Unable to marshal key as JSON: %v", err)
		}

		// Add to backend
		err = ks.store.Set(ctx, fmt.Sprintf("jwk/%s", k.ID()), payload)
		if err != nil {
			return fmt.Errorf("keystore: Unable to save key to backend: %v", err)
		}
	}

	// Synchronize the local cache
	return ks.synchronize(ctx)
}

func (ks *defaultKeyStore) Get(_ context.Context, id string) (key.Key, error) {
	ks.RLock()
	defer ks.RUnlock()

	if _, ok := ks.keys[id]; !ok {
		return nil, ErrKeyNotFound
	}

	return ks.keys[id], nil
}

func (ks *defaultKeyStore) Remove(_ context.Context, id string) error {
	ks.Lock()
	defer ks.Unlock()

	if _, ok := ks.keys[id]; ok {
		delete(ks.keys, id)
		return nil
	}

	return ErrKeyNotFound
}

func (ks *defaultKeyStore) StartMonitor(ctx context.Context) {
	// Initialize a default context
	ctx, cancel := context.WithCancel(ctx)
	ks.done = cancel

	// Fork the monitor process
	go ks.monitor(ctx)
}

func (ks *defaultKeyStore) Close() {
	if ks.done != nil {
		ks.done()
	}
}

// -----------------------------------------------------------------------------

func (ks *defaultKeyStore) monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(ks.dopts.Interval) * time.Second):
			if err := ks.synchronize(ctx); err != nil {
				log.For(ctx).Error("Unable to synchronize keystore with backend", zap.Error(err))
			}
		}
	}
}

func (ks *defaultKeyStore) synchronize(ctx context.Context) error {
	keys, err := ks.store.List(ctx, "jwk")
	if err != nil {
		return errors.Wrap(err, "keystore: Unable to synchronize keystore with backend")
	}

	// Delete keys
	for kid := range ks.keys {
		if err := ks.Remove(ctx, kid); err != nil {
			log.For(ctx).Error("Unable to remove key from local cache", zap.Error(err))
		}
	}

	// Update keys
	for _, kid := range keys {
		// Retrieve each value
		value, err := ks.store.Get(ctx, fmt.Sprintf("jwk/%s", kid))
		if err != nil {
			continue
		}

		// Decode value as Key
		holder := &keyHolder{}
		err = json.Unmarshal(value, holder)
		if err != nil {
			return err
		}

		// If key is expired ignore it
		if holder.IsExpired() {
			continue
		}

		payload := holder.Data
		if ks.dopts.Snappy {
			// Decompress buffer
			payload, err = snappy.Decode(nil, holder.Data)
			if err != nil {
				continue
			}
		}

		// Deserialize JWK
		k, err := key.FromString(payload)
		if err != nil {
			continue
		}

		// Add key to local cache
		ks.Lock()
		ks.keys[k.ID()] = k
		ks.Unlock()
	}

	return nil
}
