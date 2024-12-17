package utils

import "context"

type contextKey string

func (c contextKey) String() string {
	return "gotrue api context key " + string(c)
}

const (
	requestIDKey = contextKey("request_id")
)

func GetRequestID(ctx context.Context) string {
	obj := ctx.Value(requestIDKey)
	if obj == nil {
		return ""
	}
	return obj.(string)
}

// WithRequestID adds the provided request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}
