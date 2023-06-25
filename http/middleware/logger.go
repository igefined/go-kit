package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/igdotog/core/logger"

	"github.com/sirupsen/logrus"
)

func JsonLogger(logger *logger.Logger, excluded []string) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger, excluded})
}

type StructuredLogger struct {
	Logger      *logger.Logger
	excludedURI []string
}

type StructuredLoggerEntry struct {
	Logger     logrus.FieldLogger
	isExcluded bool
}

func isExcluded(in string, excluded []string) bool {
	for _, uri := range excluded {
		if uri == in {
			return true
		}
	}

	return false
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{
		Logger:     l.Logger,
		isExcluded: isExcluded(r.RequestURI, l.excludedURI),
	}

	if !entry.isExcluded {
		logFields := logrus.Fields{}

		logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

		if reqID := middleware.GetReqID(r.Context()); reqID != "" {
			logFields["req_id"] = reqID
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		logFields["http_scheme"] = scheme
		logFields["http_proto"] = r.Proto
		logFields["http_method"] = r.Method
		logFields["remote_addr"] = r.RemoteAddr
		logFields["user_agent"] = r.UserAgent()
		logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

		entry.Logger = entry.Logger.WithFields(logFields)
		entry.Logger.Infoln("request started")
	}

	return entry
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	if !l.isExcluded {
		l.Logger = l.Logger.WithFields(logrus.Fields{
			"resp_status": status, "resp_bytes_length": bytes,
			"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
		})

		l.Logger.Infoln("request complete")
	}
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
