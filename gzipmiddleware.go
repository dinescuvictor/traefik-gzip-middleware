package gzipmiddleware

import (
    "compress/gzip"
    "context"
    "net/http"
    "strings"
)

// Config defines the plugin's configuration.
type Config struct {
    MimeTypes []string `json:"mimeTypes,omitempty"`
}

// CreateConfig creates the default configuration.
func CreateConfig() *Config {
    return &Config{
        MimeTypes: []string{"text/html", "text/css", "application/javascript"},
    }
}

// GzipMiddleware is the main middleware struct.
type GzipMiddleware struct {
    next      http.Handler
    mimeTypes map[string]struct{}
}

// New creates a new instance of the middleware.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    mimeMap := make(map[string]struct{})
    for _, mime := range config.MimeTypes {
        mimeMap[mime] = struct{}{}
    }

    return &GzipMiddleware{
        next:      next,
        mimeTypes: mimeMap,
    }, nil
}

// ServeHTTP processes the request and compresses the response if applicable.
func (m *GzipMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    acceptEncoding := req.Header.Get("Accept-Encoding")
    if !strings.Contains(acceptEncoding, "gzip") {
        m.next.ServeHTTP(rw, req)
        return
    }

    // Enable gzip compression.
    rw.Header().Set("Content-Encoding", "gzip")
    gz := gzip.NewWriter(rw)
    defer gz.Close()

    writer := gzipResponseWriter{Writer: gz, ResponseWriter: rw}
    m.next.ServeHTTP(&writer, req)
}

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
    return w.Writer.Write(data)
}
