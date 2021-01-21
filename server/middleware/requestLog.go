package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

//HeaderRequestID gets RequestID in request header
const HeaderRequestID = "X-Request-ID"

type key int

const requestIDKey key = 0

// RequestLoggerMiddleware containing logger to log request
type RequestLoggerMiddleware struct {
	Logger *zerolog.Logger
}

// NewRequestLoggerMiddleware returns new request logger
func NewRequestLoggerMiddleware(logger *zerolog.Logger) *RequestLoggerMiddleware {
	loggerMiddleware := RequestLoggerMiddleware{
		Logger: logger,
	}
	return &loggerMiddleware
}

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get(HeaderRequestID)
	if reqID == "" {
		reqID = uuid.NewV4().String() + "-" + uuid.NewV4().String()
	}

	return context.WithValue(ctx, requestIDKey, reqID)
}

// RequestIDFromContext returns the request id from request context
func RequestIDFromContext(ctx context.Context) string {
	if ctx.Value(requestIDKey) != nil {
		return ctx.Value(requestIDKey).(string)
	}
	return ""
}

// GetMiddlewareHandler function returns middleware used to log requests
func (lm *RequestLoggerMiddleware) GetMiddlewareHandler() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := newContextWithRequestID(r.Context(), r)
		metrics := httpsnoop.CaptureMetrics(next, rw, r.WithContext(ctx))
		requestID := rw.Header().Get(HeaderRequestID)
		lm.Logger.Info().
			Str("RequestID", requestID).
			Str("Host", r.Host).
			Str("Method", r.Method).
			Str("Path", r.RequestURI).
			Str("RemoteAddr", r.RemoteAddr).
			Str("Ref", r.Referer()).
			Str("UA", r.UserAgent()).
			Int("Code", metrics.Code).
			Int("Duration", int(metrics.Duration/time.Microsecond)).
			Msg("")
	}
}
