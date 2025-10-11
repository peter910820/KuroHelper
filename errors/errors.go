package kurohelpererrors

import "errors"

var (
	// option not found error
	ErrOptionNotFound = errors.New("option: option not found")
	// option translate fail error
	ErrOptionTranslateFail = errors.New("option: value translate fail")
	// search no content for response
	ErrSearchNoContent = errors.New("search: no content for response")
	// The remote server returns a non-200 response status code
	ErrStatusCodeAbnormal = errors.New("response: server returned an error status code")
	// rate limit
	ErrRateLimit = errors.New("rate limit: rate limit, quota exhausted")
	// cache lost or expired
	ErrCacheLost = errors.New("cache: cache lost or expired")

	//ymgal invalid access token(401)
	ErrYmgalInvalidAccessToken = errors.New("ymgal: invalid access token or other 401 error")
)

// CustomID(CID) error
var (
	// cid get parameter failed
	ErrCIDGetParameterFailed = errors.New("cid: get parameter failed")
)
