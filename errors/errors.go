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
	// trying to use bangumi character list search
	ErrBangumiCharacterListSearchNotSupported = errors.New("bangumi: character list search is not currently supported")
)

// CustomID(CID) error
var (
	// cid wrong format
	ErrCIDWrongFormat = errors.New("cid: wrong format")
	// cid get parameter failed
	ErrCIDGetParameterFailed = errors.New("cid: get parameter failed")
)

// Utils error
var (
	// wrong time.Time format
	ErrTimeWrongFormat = errors.New("time: wrong format")
	// date exceeds tomorrow error
	ErrDateExceedsTomorrow = errors.New("time: date exceeds tomorrow")
)
