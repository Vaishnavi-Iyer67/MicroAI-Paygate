package main

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// bufferedWriter captures response writes in-memory so the middleware can
// decide whether to send the real response or a timeout response without
// racing with handler writes.
type bufferedWriter struct {
	buf    *bytes.Buffer
	head   http.Header
	status int
	wrote  bool
}

// newBufferedWriter returns an initialized bufferedWriter used to capture
// response headers and body from handlers without flushing to the client.
func newBufferedWriter() *bufferedWriter {
	return &bufferedWriter{
		buf:    bytes.NewBuffer(nil),
		head:   make(http.Header),
		status: http.StatusOK,
	}
}

// Header returns the local header map for the buffered response.
func (b *bufferedWriter) Header() http.Header {
	return b.head
}

func (b *bufferedWriter) Write(data []byte) (int, error) {
	b.wrote = true
	return b.buf.Write(data)
}

func (b *bufferedWriter) WriteString(s string) (int, error) {
	b.wrote = true
	return b.buf.WriteString(s)
}

// WriteHeader captures the status code but does not flush to the client.
func (b *bufferedWriter) WriteHeader(statusCode int) {
	b.status = statusCode
}

// WriteHeaderNow is a no-op to avoid flushing headers to the client while
// we're buffering.
func (b *bufferedWriter) WriteHeaderNow() {
	// Intentionally left blank - avoid flushing
}

// Status returns the status code that the handler set (or 200 by default).
func (b *bufferedWriter) Status() int {
	if b.status == 0 {
		return http.StatusOK
	}
	return b.status
}

// flushTo writes buffered headers and body to the real writer.
func (b *bufferedWriter) flushTo(w http.ResponseWriter) {
	for k, vv := range b.head {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(b.Status())
	_, _ = w.Write(b.buf.Bytes())
}

// RequestTimeoutMiddleware applies a context timeout to the request and
// buffers handler output. If the context deadline is exceeded, the middleware
// returns 504 and discards the handler response. This avoids concurrent
// response writes and ensures safe behavior with Gin.
func RequestTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protect against non-positive timeouts by treating them as zero
		// which results in immediate deadline behavior if <= 0. For safety we
		// allow zero (caller can disable middleware by not adding it).
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		origWriter := c.Writer
		bw := newBufferedWriter()
		// replace the gin writer with a shim that uses bw
		c.Writer = &responseWriterShim{bw: bw}

		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					panicChan <- r
				}
			}()
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// Handler finished before deadline: flush buffered response
			c.Writer = origWriter
			bw.flushTo(origWriter)
			return
		case p := <-panicChan:
			// Re-panic in this goroutine so upstream recovery runs as expected
			panic(p)
		case <-ctx.Done():
			// Timeout exceeded — send 504 using the original writer and
			// discard buffered handler response.
			c.Writer = origWriter
			origWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
			origWriter.WriteHeader(504)
			_, _ = origWriter.Write([]byte(`{"error":"Gateway Timeout","message":"Request exceeded maximum allowed time"}`))
			return
		}
	}
}

// responseWriterShim adapts bufferedWriter to satisfy gin.ResponseWriter so
// handlers that call c.Writer/SetHeader interact with the buffered headers
// and body. It forwards writes to the underlying bufferedWriter instance.
type responseWriterShim struct {
	bw *bufferedWriter
}

func (rws *responseWriterShim) Header() http.Header { return rws.bw.Header() }
func (rws *responseWriterShim) Write(data []byte) (int, error) { return rws.bw.Write(data) }
func (rws *responseWriterShim) WriteString(s string) (int, error) { return rws.bw.WriteString(s) }
func (rws *responseWriterShim) WriteHeader(statusCode int) { rws.bw.WriteHeader(statusCode) }
func (rws *responseWriterShim) WriteHeaderNow() { rws.bw.WriteHeaderNow() }
func (rws *responseWriterShim) Status() int { return rws.bw.Status() }
func (rws *responseWriterShim) Written() bool { return rws.bw.wrote }
func (rws *responseWriterShim) Size() int { return rws.bw.buf.Len() }
func (rws *responseWriterShim) WriteHeaderNowWithoutLock() {}


