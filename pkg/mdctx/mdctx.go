package mdctx

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/pkg/util"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// New returns a context with a random operation ID. The context should be used throughout single request processing to correlate all logs
// and actions using its operation ID.
func New() context.Context {
	ctx := context.Background()
	return withValue(ctx, operationIdKey, util.RandomAlphanumericString(operationIdLength))
}

// WithCorrelationId returns a copy of context with added correlation ID. Unlike operation ID, correlation ID should come from clients'
// requests and may be associated with multiple requests.
func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return withValue(ctx, correlationIdKey, correlationId)
}

// WithRequestMethod returns a copy of context with added HTTP method (e.g. GET).
func WithRequestMethod(ctx context.Context, requestMethod string) context.Context {
	return withValue(ctx, requestMethodKey, requestMethod)
}

// WithRequestUri returns a copy of context with added URI (e.g. /api/customers?pageSize=10).
func WithRequestUri(ctx context.Context, uri string) context.Context {
	return withValue(ctx, requestUriKey, uri)
}

// WithClientIp returns a copy of context with added client IP (e.g. 80.77.213.103).
func WithClientIp(ctx context.Context, clientIp string) context.Context {
	return withValue(ctx, clientIpKey, clientIp)
}

func withValue(ctx context.Context, key mdcKey, value string) context.Context {
	return context.WithValue(ctx, key, mdcValue(value))
}

// CorrelationId extracts correlation ID from context to use when communicating with other services. Correlation ID is returned if it's set.
// Otherwise, operation ID is returned if it's set. If neither is set then empty string is returned.
func CorrelationId(ctx context.Context) string {
	if value, ok := ctx.Value(correlationIdKey).(mdcValue); ok {
		return string(value)
	}
	if value, ok := ctx.Value(operationIdKey).(mdcValue); ok {
		return string(value)
	}
	return ""
}

// OperationId extracts operation ID from context. Empty string is returned if it's not set.
func OperationId(ctx context.Context) string {
	if value, ok := ctx.Value(operationIdKey).(mdcValue); ok {
		return string(value)
	}
	return ""
}

// Debugf logs at debug level. Arguments are handled in the manner of fmt.Printf.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logf(logLevelDebug, ctx, format, args...)
}

// Infof logs at info level. Arguments are handled in the manner of fmt.Printf.
func Infof(ctx context.Context, format string, args ...interface{}) {
	logf(logLevelInfo, ctx, format, args...)
}

// Warnf logs at warn level. Arguments are handled in the manner of fmt.Printf.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logf(logLevelWarn, ctx, format, args...)
}

// Errorf logs at error level. Arguments are handled in the manner of fmt.Printf.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logf(logLevelError, ctx, format, args...)
}

// Fatalf logs at fatal level and exists with non-zero exit code. Arguments are handled in the manner of fmt.Printf.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logf(logLevelFatal, ctx, format, args...)
	os.Exit(1)
}

func logf(logLevel logLevel, ctx context.Context, format string, args ...interface{}) {
	if ctx == nil {
		ctx = context.TODO()
	}

	decoratedFormat := fmt.Sprintf("%s %s %s]   %s   ]%s", timestamp(), logLevel, loggingFileAndLine(), format, mdcLabels(ctx))
	logger.Printf(decoratedFormat, args...)
}

func timestamp() string {
	return currentTimeHook().UTC().Format("2006-01-02T15:04:05.000Z")
}

// For mocking in unit tests.
var currentTimeHook = time.Now

func loggingFileAndLine() string {
	_, filepath, lineNumber, ok := runtime.Caller(3)
	if !ok {
		return "???"
	}

	pathSegments := strings.Split(filepath, "/")
	file := pathSegments[len(pathSegments)-1]
	// Trim the file to maintain coherent logging format (20 corresponds to %20s in Sprintf below, which pads filename to 20 characters)
	if len(file) > 20 {
		file = "..." + file[len(file)-17:]
	}

	return fmt.Sprintf("%20s:%-4d", file, lineNumber)
}

func mdcLabels(ctx context.Context) string {
	var mdcString string
	for _, key := range mdcKeys {
		if value, ok := ctx.Value(key).(mdcValue); ok {
			mdcString += fmt.Sprintf("[%s=%s]", key, value)
		}
	}
	return mdcString
}

type (
	logLevel string
	mdcKey   string
	mdcValue string
)

const (
	operationIdLength        = 10
	operationIdKey    mdcKey = "operation-id"
	correlationIdKey  mdcKey = "correlation-id"
	requestMethodKey  mdcKey = "http"
	requestUriKey     mdcKey = "uri"
	clientIpKey       mdcKey = "client-ip"
	// All log levels should have the same length
	logLevelDebug logLevel = "DEBUG"
	logLevelInfo  logLevel = "INFO "
	logLevelWarn  logLevel = "WARN "
	logLevelError logLevel = "ERROR"
	logLevelFatal logLevel = "FATAL"
)

var (
	// mdcKeys contains all defined mdc keys and defines the order in which they appear in log messages.
	mdcKeys = []mdcKey{
		operationIdKey,
		correlationIdKey,
		requestMethodKey,
		requestUriKey,
		clientIpKey,
	}
	logger = log.New(os.Stderr, "", 0)
)
