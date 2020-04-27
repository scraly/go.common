package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

func TestLifecycle_cancel_err_empty(t *testing.T) {
	gomega.RegisterTestingT(t)

	cancel = nil
	gomega.Expect(cancel).Should(gomega.BeNil())

	ctx := context.Background()

	// cancel
	err := Cancel(ctx)

	gomega.Expect(err).Should(gomega.MatchError("Look like this context aren't managed by this package. No CancelFunc found"))
}

func TestLifecycle_create(t *testing.T) {
	gomega.RegisterTestingT(t)

	ctx := Create(context.Background())

	gomega.Expect(ctx).ShouldNot(gomega.BeNil())

	ctxID := ctx.Value(ctxKey).(string)
	gomega.Expect(ctxID).ShouldNot(gomega.BeNil())
	gomega.Expect(ctxID).ShouldNot(gomega.BeEmpty())

	gomega.Expect(cancel).ShouldNot(gomega.BeNil())
	gomega.Expect(cancel[ctxID]).ShouldNot(gomega.BeNil())
}

func TestLifecycle_cancel_ok(t *testing.T) {
	gomega.RegisterTestingT(t)

	ctx := Create(context.Background())

	gomega.Expect(ctx).ShouldNot(gomega.BeNil())

	ch := make(chan bool)

	//go func
	go func() {
		select {
		case <-ctx.Done():
			ch <- true
		case <-time.After(1 * time.Second):
			ch <- false
		}
	}()

	// cancel
	err := Cancel(ctx)
	gomega.Expect(err).Should(gomega.BeNil())

	result := <-ch

	gomega.Expect(result).Should(gomega.BeTrue())
}

func TestLifecycle_cancel_err_ctx_not_managed(t *testing.T) {
	gomega.RegisterTestingT(t)

	_ = Create(context.Background())
	gomega.Expect(cancel).ShouldNot(gomega.BeNil())

	ctx := context.Background()

	// cancel
	err := Cancel(ctx)

	gomega.Expect(err).Should(gomega.MatchError("Look like this context aren't managed by this package. No CancelFunc found"))
}

func TestLifecycle_cancel_err_missing_func(t *testing.T) {
	gomega.RegisterTestingT(t)

	_ = Create(context.Background())
	gomega.Expect(cancel).ShouldNot(gomega.BeNil())

	ctx := context.WithValue(context.Background(), ctxKey, "blabla")

	// cancel
	err := Cancel(ctx)

	gomega.Expect(err).Should(gomega.MatchError("CancelFunc missing. Unable to make a call on it"))
}
