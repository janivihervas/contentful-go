package contentful

import "go.opencensus.io/trace"

func addSpanError(span *trace.Span, statusCode int32, err error) {
	span.SetStatus(trace.Status{Code: statusCode, Message: err.Error()})
	// For Jaeger
	span.AddAttributes(trace.BoolAttribute("error", true))
}
