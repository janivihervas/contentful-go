package contentful

import "go.opencensus.io/trace"

func spanError(span *trace.Span, statusCode int32, err error) {
	span.SetStatus(trace.Status{Code: statusCode, Message: err.Error()})
	span.AddAttributes(trace.BoolAttribute("error", true))
}
