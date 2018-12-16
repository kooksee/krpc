// Commons for HTTP handling
package krpcs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/net/netutil"

	types "github.com/kooksee/krpc/types"
)

// Config is an RPC server configuration.
type Config struct {
	MaxOpenConnections int
}

const (
	// maxBodyBytes controls the maximum number of bytes the
	// server will read parsing the request body.
	maxBodyBytes = int64(1024 * 1024) // 1MB
)

// StartHTTPServer starts an HTTP server on listenAddr with the given handler.
// It wraps handler with RecoverAndLogHandler.
func StartHTTPServer(listenAddr string, handler http.Handler, config Config) net.Listener {
	parts := strings.SplitN(listenAddr, "://", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)", listenAddr))
	}

	proto, addr := parts[0], parts[1]

	log.Info().Msgf("Starting RPC HTTP server on %s", listenAddr)
	listener, err := net.Listen(proto, addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on %v: %v", listenAddr, err))
	}

	if config.MaxOpenConnections > 0 {
		listener = netutil.LimitListener(listener, config.MaxOpenConnections)
	}

	if err := http.Serve(listener, RecoverAndLogHandler(maxBytesHandler{h: handler, n: maxBodyBytes})); err != nil {
		log.Error().Err(err).Msg("RPC HTTP server stopped")
		panic(err.Error())
	}

	return listener
}

// StartHTTPAndTLSServer starts an HTTPS server on listenAddr with the given
// handler.
// It wraps handler with RecoverAndLogHandler.
func StartHTTPAndTLSServer(listenAddr string, handler http.Handler, certFile, keyFile string, config Config) {
	var proto, addr string
	parts := strings.SplitN(listenAddr, "://", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)", listenAddr))
	}
	proto, addr = parts[0], parts[1]

	log.Info().Msgf("Starting RPC HTTPS server on %s (cert: %q, key: %q)", listenAddr, certFile, keyFile)

	listener, err := net.Listen(proto, addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on %v: %v", listenAddr, err))
	}

	if config.MaxOpenConnections > 0 {
		listener = netutil.LimitListener(listener, config.MaxOpenConnections)
	}

	if err := http.ServeTLS(listener, RecoverAndLogHandler(maxBytesHandler{h: handler, n: maxBodyBytes}), certFile, keyFile); err != nil {
		log.Error().Err(err).Msg("RPC HTTPS server stopped")
		panic(err.Error())
	}
}

func WriteRPCResponseHTTPError(w http.ResponseWriter, httpCode int, res types.RPCResponse) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(jsonBytes) // nolint: errcheck, gas
}

func WriteRPCResponseHTTP(w http.ResponseWriter, res types.RPCResponse) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonBytes) // nolint: errcheck, gas
}

//-----------------------------------------------------------------------------

// Wraps an HTTP handler, adding error logging.
// If the inner function panics, the outer function recovers, logs, sends an
// HTTP 500 error response.
func RecoverAndLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to remember the status
		rww := &ResponseWriterWrapper{-1, w}
		begin := time.Now()

		// Common headers
		origin := r.Header.Get("Origin")
		rww.Header().Set("Access-Control-Allow-Origin", origin)
		rww.Header().Set("Access-Control-Allow-Credentials", "true")
		rww.Header().Set("Access-Control-Expose-Headers", "X-Server-Time")
		rww.Header().Set("X-Server-Time", fmt.Sprintf("%v", begin.Unix()))

		defer func() {
			// Send a 500 error if a panic happens during a handler.
			// Without this, Chrome & Firefox were retrying aborted ajax requests,
			// at least to my localhost.
			if e := recover(); e != nil {

				// If RPCResponse
				if res, ok := e.(types.RPCResponse); ok {
					WriteRPCResponseHTTP(rww, res)
				} else {
					// For the rest,
					log.Error().Str("stack", string(debug.Stack())).Msg("Panic in RPC HTTP handler")
					rww.WriteHeader(http.StatusInternalServerError)
					WriteRPCResponseHTTP(rww, types.RPCInternalError("", e.(error)))
				}
			}

			// Finally, log.
			durationMS := time.Since(begin).Nanoseconds() / 1000000
			if rww.Status == -1 {
				rww.Status = 200
			}
			log.Info().Str("method", r.Method).Str("url", r.URL.String()).Int("status", rww.Status).Int64("duration", durationMS).Str("remoteAddr", r.RemoteAddr).Msg("Served RPC HTTP response")
		}()

		handler.ServeHTTP(rww, r)
	})
}

// Remember the status for logging
type ResponseWriterWrapper struct {
	Status int
	http.ResponseWriter
}

func (w *ResponseWriterWrapper) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// implements http.Hijacker
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (h maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.n)
	h.h.ServeHTTP(w, r)
}
