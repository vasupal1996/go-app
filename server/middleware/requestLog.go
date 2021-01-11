package middleware

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
)

//HeaderRequestID gets RequestID in request header
const HeaderRequestID = "X-Request-ID"

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

// GetMiddlewareHandler function returns middleware used to log requests
func (lm *RequestLoggerMiddleware) GetMiddlewareHandler() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		metrics := httpsnoop.CaptureMetrics(next, rw, r)
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
