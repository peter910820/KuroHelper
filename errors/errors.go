package errors

import "errors"

var (
	ErrOptionNotFound      = errors.New("[GetOptionsErr]: option not found")
	ErrOptionTranslateFail = errors.New("[GetOptionsErr]: value translate fail")
)
