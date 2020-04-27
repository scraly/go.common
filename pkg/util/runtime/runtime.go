package runtime

import (
	"context"
	"fmt"

	"github.com/scraly/go.common/pkg/util/contexthelper"
	"github.com/nats-io/nuid"
)

var cancel map[string]context.CancelFunc
var ctxKey = contexthelper.DefaultContextKey("RunID")

// Create will create a Context and remind the context.CancelFunc
func Create(parent context.Context) context.Context {
	ctx, cancelFunc := context.WithCancel(parent)
	ctxID := nuid.Next()
	ctx = context.WithValue(ctx, ctxKey, ctxID)

	if cancel == nil {
		cancel = make(map[string]context.CancelFunc)
	}
	cancel[ctxID] = cancelFunc

	return ctx
}

// Cancel will call the previous registered context.CancelFunc
func Cancel(ctx context.Context) error {
	if cancel != nil && ctx.Value(ctxKey) != nil {
		ctxID := ctx.Value(ctxKey).(string)
		if cancel[ctxID] != nil {
			cancel[ctxID]()
			return nil
		}
		return fmt.Errorf("CancelFunc missing. Unable to make a call on it")
	}
	return fmt.Errorf("Look like this context aren't managed by this package. No CancelFunc found")
}
