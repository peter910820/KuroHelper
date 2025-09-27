package kurohelpererrors

import "errors"

var (
	ErrOptionNotFound      = errors.New("option: option not found")
	ErrOptionTranslateFail = errors.New("option: value translate fail")

	ErrSearchNoContent = errors.New("search: no content for search")

	// The remote server returns a non-200 response status code
	ErrStatusCodeAbnormal = errors.New("server returned an error status code")
	// rate limit
	ErrRateLimit = errors.New("rate limit, quota exhausted")
)
