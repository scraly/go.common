package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/scraly/go.common/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	// DefaultRotation delay for token rotation
	DefaultRotation = 60 * time.Minute
)

// RotatingTokenSource declare token source rotatable contract
type RotatingTokenSource interface {
	oauth2.TokenSource
	Client(context.Context) *http.Client
	Start(context.Context)
}

// RotatingTokenSource represents OAuth 2.0 token sourcing implementation
type rotatingTokenSource struct {
	src         oauth2.TokenSource
	token       *oauth2.Token
	mu          sync.RWMutex
	lastRefresh time.Time
}

// NewRotatingTokenSource returns a TokenSource implementation for Entry server
func NewRotatingTokenSource(ctx context.Context, src oauth2.TokenSource, delay uint16) (RotatingTokenSource, error) {
	ts := &rotatingTokenSource{
		src: src,
	}
	DefaultRotation = time.Duration(delay) * time.Minute

	// Rotate on start
	if err := ts.getToken(ctx); err != nil {
		return nil, errors.Wrap(err, "unable to retrieve initial access token")
	}

	return ts, nil
}

// -----------------------------------------------------------------------------

// Token retrieve an access token from provider
func (ts *rotatingTokenSource) Token() (*oauth2.Token, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.token, nil
}

// Retrieve a client with rotating tokensource
func (ts *rotatingTokenSource) Client(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, ts)
}

func (ts *rotatingTokenSource) Start(ctx context.Context) {
	// Start rotating routine
	go ts.tokenRenewalLoop(ctx)
}

// -----------------------------------------------------------------------------

func (ts *rotatingTokenSource) tokenRenewalLoop(ctx context.Context) {

	ticker := time.NewTicker(DefaultRotation)

	for {
		select {
		case <-ctx.Done():
			// Exiting rotator loop
			ticker.Stop()
			break
		case <-ticker.C:
			// After expiration delay do the rotation
			log.For(ctx).Debug("Starting token rotation")
			if err := ts.getToken(ctx); err != nil {
				break
			}
		}
	}

}

func (ts *rotatingTokenSource) getToken(ctx context.Context) error {

	// Check refresh time
	if ts.lastRefresh.Add(DefaultRotation).After(time.Now()) {
		// Skip rotation
		log.For(ctx).Debug("Rotation skipped")
		return nil
	}

	// Refresh token from original tokensource
	token, err := ts.src.Token()

	if err != nil {
		return errors.Wrap(err, "unable to retrieve access token from original token source")
	}
	newToken := sha256.Sum256([]byte("Bearer " + token.AccessToken))
	log.For(ctx).Debug(base64.StdEncoding.EncodeToString(newToken[:]))
	// Assign new token
	ts.mu.Lock()
	ts.token = token
	ts.lastRefresh = time.Now()
	ts.mu.Unlock()

	// Return no error
	return nil
}
