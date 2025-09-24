package errors

import "errors"

var (
	ErrOptionNotFound      = errors.New("option: option not found")
	ErrOptionTranslateFail = errors.New("option: value translate fail")

	ErrVndbNoResult = errors.New("vndb: no result for response")

	ErrSearchNoContent = errors.New("search: no content for search")
)
