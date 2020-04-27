package session

type contextKey string

func (c contextKey) String() string {
	return "github.com/scraly/go.common/pkg/web/middleware/" + string(c)
}
