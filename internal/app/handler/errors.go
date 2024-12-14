package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tranminhquanq/gomess/internal/utils"
)

type HTTPError struct {
	HTTPStatus      int    `json:"code"`                 // do not rename the JSON tags!
	ErrorCode       string `json:"error_code,omitempty"` // do not rename the JSON tags!
	Message         string `json:"message"`              // do not rename the JSON tags!
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
	ErrorID         string `json:"error_id,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.HTTPStatus, e.Message)
}

func (e *HTTPError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

type HTTPErrorResponse20240101 struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func HandleResponseError(err error, w http.ResponseWriter, r *http.Request) {
	errorID := utils.GetRequestID(r.Context())

	switch e := err.(type) {
	case *HTTPError:
		switch {
		case e.HTTPStatus >= http.StatusInternalServerError:
			e.ErrorID = errorID
			// this will get us the stack trace too
			logrus.Error(r.Context(), e)
		case e.HTTPStatus == http.StatusTooManyRequests:
			logrus.WithError(e.Cause()).Warn(e.Error())
		default:
			logrus.WithError(e.Cause()).Info(e.Error())
		}

	case ErrorCause:
		HandleResponseError(e.Cause(), w, r)

	default:
		resp := HTTPErrorResponse20240101{
			Code:    ErrorCodeUnexpectedFailure,
			Message: "Unexpected failure, please check server logs for more information",
		}
		if jsonErr := sendJSON(w, http.StatusInternalServerError, resp); jsonErr != nil && jsonErr != context.DeadlineExceeded {
			logrus.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
		}
	}
}

func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				// logEntry := observability.GetLogEntry(r)
				// if logEntry != nil {
				// 	logEntry.Panic(rvr, debug.Stack())
				// } else {
				// 	fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				// 	debug.PrintStack()
				// }

				se := &HTTPError{
					HTTPStatus: http.StatusInternalServerError,
					Message:    http.StatusText(http.StatusInternalServerError),
				}
				HandleResponseError(se, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// ErrorCause is an error interface that contains the method Cause() for returning root cause errors
type ErrorCause interface {
	Cause() error
}

// func generateFrequencyLimitErrorMessage(timeStamp *time.Time, maxFrequency time.Duration) string {
// 	now := time.Now()
// 	left := timeStamp.Add(maxFrequency).Sub(now) / time.Second
// 	return fmt.Sprintf("For security purposes, you can only request this after %d seconds.", left)
// }
