package utils

import (
	kurohelpererrors "kurohelper/errors"
	"os"
	"strconv"
	"time"
)

// check if the string is English
func IsEnglish(r rune) bool {
	if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
		return false
	}
	return true
}

// 取得env(轉型int版)
func GetEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}
	return def
}

func ParseYYYYMMDD(s string) (time.Time, error) {
	if len(s) != 8 {
		return time.Time{}, kurohelpererrors.ErrTimeWrongFormat
	}

	if _, err := strconv.Atoi(s); err != nil {
		return time.Time{}, kurohelpererrors.ErrTimeWrongFormat
	}

	t, err := time.Parse("20060102", s)
	if err != nil {
		return time.Time{}, kurohelpererrors.ErrTimeWrongFormat
	}

	return t, nil
}
