package session

import (
	"context"
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

const (
	contextSessions = contextKey("sessions")
)

const (
	// FlashErr is the key for errors
	FlashErr = "_flash_err"
	// FlashWarn is the key for warnings
	FlashWarn = "_flash_warn"
	// FlashInfo is the key for informations
	FlashInfo = "_flash_info"
)

// FromContext returns renderer from context object
func FromContext(ctx context.Context) (*sessions.Session, bool) {
	value, ok := ctx.Value(contextSessions).(*sessions.Session)
	return value, ok
}

// NewMiddleware is used to expose render template to request
func NewMiddleware(name string, store sessions.Store) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, name)
			defer gcontext.Clear(r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, contextSessions, session)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
