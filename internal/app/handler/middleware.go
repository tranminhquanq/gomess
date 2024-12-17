package handler

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// timeoutResponseWriter is a http.ResponseWriter that queues up a response
// body to be sent if the serving completes before the context has exceeded its
// deadline.
type timeoutResponseWriter struct {
	sync.Mutex

	header      http.Header
	wroteHeader bool
	snapHeader  http.Header // snapshot of the header at the time WriteHeader was called
	statusCode  int
	buf         bytes.Buffer
}

func (t *timeoutResponseWriter) Header() http.Header {
	t.Lock()
	defer t.Unlock()

	return t.header
}

func (t *timeoutResponseWriter) Write(bytes []byte) (int, error) {
	t.Lock()
	defer t.Unlock()

	if !t.wroteHeader {
		t.writeHeaderLocked(http.StatusOK)
	}

	return t.buf.Write(bytes)
}

func (t *timeoutResponseWriter) WriteHeader(statusCode int) {
	t.Lock()
	defer t.Unlock()

	t.writeHeaderLocked(statusCode)
}

func (t *timeoutResponseWriter) writeHeaderLocked(statusCode int) {
	if t.wroteHeader {
		// ignore multiple calls to WriteHeader
		// once WriteHeader has been called once, a snapshot of the header map is taken
		// and saved in snapHeader to be used in finallyWrite
		return
	}

	t.statusCode = statusCode
	t.wroteHeader = true
	t.snapHeader = t.header.Clone()
}

func (t *timeoutResponseWriter) finallyWrite(w http.ResponseWriter) {
	t.Lock()
	defer t.Unlock()

	dst := w.Header()
	for k, vv := range t.snapHeader {
		dst[k] = vv
	}

	if !t.wroteHeader {
		t.statusCode = http.StatusOK
	}

	w.WriteHeader(t.statusCode)
	if _, err := w.Write(t.buf.Bytes()); err != nil {
		logrus.WithError(err).Warn("Write failed")
	}
}

func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			timeoutWriter := &timeoutResponseWriter{
				header: make(http.Header),
			}

			panicChan := make(chan any, 1)
			serverDone := make(chan struct{})
			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()

				next.ServeHTTP(timeoutWriter, r.WithContext(ctx))
				close(serverDone)
			}()

			select {
			case p := <-panicChan:
				panic(p)

			case <-serverDone:
				timeoutWriter.finallyWrite(w)

			case <-ctx.Done():
				err := ctx.Err()

				if err == context.DeadlineExceeded {
					httpError := &HTTPError{
						HTTPStatus: http.StatusGatewayTimeout,
						ErrorCode:  ErrorCodeRequestTimeout,
						Message:    "Processing this request timed out, please retry after a moment.",
					}

					httpError = httpError.WithInternalError(err)

					HandleResponseError(httpError, w, r)
				} else {
					// unrecognized context error, so we should wait for the server to finish
					// and write out the response
					<-serverDone

					timeoutWriter.finallyWrite(w)
				}
			}
		})
	}
}
