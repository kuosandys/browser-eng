package errors

import "errors"

var (
	ErrUnsupportedURLScheme  = errors.New("unsupported URL scheme")
	ErrUnsupportedMediaType  = errors.New("unsupported media type")
	ErrMissingLocationHeader = errors.New("missing Location header")
	ErrMaxRedirectsReached   = errors.New("max redirects reached")
)
