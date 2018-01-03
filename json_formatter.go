package epiclogger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/facebookgo/stack"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	errorReporting "google.golang.org/api/clouderrorreporting/v1beta1"
	logging "google.golang.org/api/logging/v2beta1"
	"google.golang.org/grpc/metadata"
)

// EpicFormatter is similar to logrus.JSONFormatter but with log level that are recongnized
// by kubernetes fluentd.
type EpicFormatter struct{}

func isError(entry *log.Entry) bool {
	if entry != nil {
		switch entry.Level {
		case log.ErrorLevel:
			return true
		case log.FatalLevel:
			return true
		case log.PanicLevel:
			return true
		}
	}
	return false
}

func retrieveMetaData(ctx context.Context) (data map[string][]string) {
	var ok bool
	data, ok = metadata.FromContext(ctx)
	if !ok {
		fmt.Println("Failed to retrieve metadata")
	}
	return
}

func contains(item string, array []string) bool {
	for _, v := range array {
		if item == v {
			return true
		}
	}
	return false
}

func prefixFieldClashes(data log.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["message"]; ok {
		data["fields.message"] = m
	}

	if l, ok := data["severity"]; ok {
		data["fields.severity"] = l
	}
}

func getSeverity(level log.Level) string {
	switch level {
	case log.FatalLevel:
		return "CRITICAL"
	case log.PanicLevel:
		return "CRITICAL"
	default:
		return strings.ToUpper(level.String())
	}
}

// Format the log entry. Implements logrus.Formatter.
func (f *EpicFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+3)
	var httpReq *logging.HttpRequest
	for k, v := range entry.Data {
		switch x := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = x.Error()
		case *http.Request:
			httpReq = &logging.HttpRequest{
				Referer:       x.Referer(),
				RemoteIp:      x.RemoteAddr,
				RequestMethod: x.Method,
				RequestUrl:    x.URL.String(),
				UserAgent:     x.UserAgent(),
			}

		case *logging.HttpRequest:
			httpReq = x

		case context.Context:
			metaData := retrieveMetaData(x)
			if authorID, ok := metaData["author_id"]; ok {
				data["userId"] = authorID[0]
			}

			if authorName, ok := metaData["author_name"]; ok {
				data["user"] = authorName[0]
			}

			if correlationID, ok := metaData["correlation_id"]; ok {
				data["correlationId"] = correlationID[0]
			}
			for key, value := range grpc_ctxtags.Extract(x).Values() {
				data[key] = fmt.Sprintf("%v", value)
			}

		default:
			data[k] = v
		}
	}

	if data["grpc.method"] != nil {
		httpReq = &logging.HttpRequest{
			RequestMethod: "POST",
			RequestUrl:    data["grpc.method"].(string),
		}
	}

	prefixFieldClashes(data)
	payload := preparePayload(entry, data, httpReq)
	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

func preparePayload(entry *log.Entry, data log.Fields, httpReq *logging.HttpRequest) map[string]interface{} {
	data["time"] = entry.Time.Format(time.RFC3339)
	data["message"] = entry.Message
	data["severity"] = getSeverity(entry.Level)
	// The error reporting payload JSON schema is defined in:
	// https://cloud.google.com/error-reporting/docs/formatting-error-messages
	// Which reflects the structure of the ErrorEvent type in:
	// https://godoc.org/google.golang.org/api/clouderrorreporting/v1beta1
	if isError(entry) {
		errorEvent := buildErrorReportingEvent(entry, data, httpReq)
		errorStructPayload, err := json.Marshal(errorEvent)
		if err != nil {
			log.Printf("error marshaling error reporting data: %s", err.Error())
		}
		var errorJSONPayload map[string]interface{}
		err = json.Unmarshal(errorStructPayload, &errorJSONPayload)
		if err != nil {
			log.Printf("error parsing error reporting data: %s", err.Error())
		}
		for k, v := range data {
			if !contains(k, []string{"service", "version", "caller", "user", "stack", "message"}) {
				errorJSONPayload[k] = v
			}
		}
		return errorJSONPayload
	}
	if httpReq != nil {
		data["httpRequest"] = httpReq
	}
	return data
}

func buildErrorReportingEvent(entry *log.Entry, data log.Fields, httpReq *logging.HttpRequest) errorReporting.ReportedErrorEvent {
	errorEvent := errorReporting.ReportedErrorEvent{
		EventTime:      entry.Time.Format(time.RFC3339),
		Message:        entry.Message,
		ServiceContext: &errorReporting.ServiceContext{},
		Context:        &errorReporting.ErrorContext{},
	}
	if data["service"] != nil {
		errorEvent.ServiceContext.Service = data["service"].(string)
	}
	if data["version"] != nil {
		errorEvent.ServiceContext.Version = data["version"].(string)
	}
	if data["stack"] != nil {
		errorEvent.Message += data["stack"].(stack.Stack).String()
	}

	if data["user"] != nil {
		errorEvent.Context.User = data["user"].(string)
	}
	// Assumes that caller stack frame information of type
	// github.com/facebookgo/stack.Frame has been added.
	// Possibly via a library like github.com/Gurpartap/logrus-stack
	if entry.Data["caller"] != nil {
		caller := entry.Data["caller"].(stack.Frame)
		errorEvent.Context.ReportLocation = &errorReporting.SourceLocation{
			FilePath:     caller.File,
			FunctionName: caller.Name,
			LineNumber:   int64(caller.Line),
		}
	}
	if httpReq != nil {
		errRepHTTPRequest := &errorReporting.HttpRequestContext{
			Method:    httpReq.RequestMethod,
			Referrer:  httpReq.Referer,
			RemoteIp:  httpReq.RemoteIp,
			Url:       httpReq.RequestUrl,
			UserAgent: httpReq.UserAgent,
		}
		errorEvent.Context.HttpRequest = errRepHTTPRequest
	}
	return errorEvent
}
