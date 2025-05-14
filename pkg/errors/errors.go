package errors

import (
	"log"
	"net/http"
	"os"
)

type LocalizeMessage struct {
	Key  string                 `json:"key"`
	Vars map[string]interface{} `json:"vars,omitempty"`
}

type Extension struct {
	Detail        interface{}       `json:"detail,omitempty"`
	Message       string            `json:"message,omitempty"`
	LocaleMessage *LocalizeMessage  `json:"locale_message,omitempty"`
	Scope         string            `json:"scope,omitempty"`
	Location      string            `json:"location,omitempty"`
	ErrorCode     string            `json:"error_code,omitempty"`
	StatusCode    int               `json:"status_code"`
	Extra         map[string]string `json:"-"`
}

func (e *Extension) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.LocaleMessage != nil && e.LocaleMessage.Key != "" {
		return e.LocaleMessage.Key
	}
	return http.StatusText(e.StatusCode)
}

type Option func(*Extension)

func WithDetail(detail interface{}) Option {
	return func(e *Extension) { e.Detail = detail }
}

func WithMessage(msg string) Option {
	return func(e *Extension) { e.Message = msg }
}

func WithScope(scope string) Option {
	return func(e *Extension) { e.Scope = scope }
}

func WithLocation(loc string) Option {
	return func(e *Extension) { e.Location = loc }
}

func WithErrorCode(code string) Option {
	return func(e *Extension) { e.ErrorCode = code }
}

func WithLocalizedMsg(key string, vars map[string]interface{}) Option {
	return func(e *Extension) {
		e.LocaleMessage = &LocalizeMessage{Key: key, Vars: vars}
	}
}

func WithExtra(key, value string) Option {
	return func(e *Extension) {
		if e.Extra == nil {
			e.Extra = make(map[string]string)
		}
		e.Extra[key] = value
	}
}

func BadRequest(opts ...Option) error {
	e := &Extension{StatusCode: http.StatusBadRequest}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func Unauthorized(opts ...Option) error {
	e := &Extension{StatusCode: http.StatusUnauthorized}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func NotFound(opts ...Option) error {
	e := &Extension{StatusCode: http.StatusNotFound}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func InternalServerError(opts ...Option) error {
	e := &Extension{StatusCode: http.StatusInternalServerError}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func TooManyRequests(opts ...Option) error {
	e := &Extension{StatusCode: http.StatusTooManyRequests}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func LogAndPanic(err error) {
	if extErr, ok := err.(*Extension); ok {
		log.Printf("[FATAL] %s/%s - %s\n%s",
			extErr.Scope,
			extErr.Location,
			extErr.Message,
			extErr.Detail,
		)
	} else {
		log.Printf("[FATAL] Unknown error: %+v", err)
	}
	os.Exit(1)
}
