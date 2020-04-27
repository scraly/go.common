package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/scraly/go.common/pkg/log"
)

type gzipResponseWriter struct {
	*gzip.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

// Handler is the Gzip request/response middleware
func Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		gw := gzip.NewWriter(w)
		defer log.SafeClose(gw, "Error while closing Gzip writer")

		w = &gzipResponseWriter{gw, w}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
