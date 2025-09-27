package vndb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	kurohelpererrors "kurohelper/errors"
)

type rateLimitStruct struct {
	Quota     int
	ResetTime time.Time
	RWMu      sync.RWMutex
}

var (
	rateLimitRecord = rateLimitStruct{
		Quota:     40,
		ResetTime: time.Now().Add(1 * time.Minute),
	}
)

func sendGetRequest(apiRoute string) ([]byte, error) {
	if !rateLimit(1) {
		return nil, kurohelpererrors.ErrRateLimit
	}

	resp, err := http.Get(os.Getenv("VNDB_ENDPOINT") + apiRoute)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func sendPostRequest(apiRoute string, jsonBytes []byte) ([]byte, error) {
	if !rateLimit(1) {
		return nil, kurohelpererrors.ErrRateLimit
	}

	resp, err := http.Post(os.Getenv("VNDB_ENDPOINT")+apiRoute, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func rateLimit(quota int) bool {
	rateLimitRecord.RWMu.Lock()
	defer rateLimitRecord.RWMu.Unlock()

	now := time.Now()
	if now.After(rateLimitRecord.ResetTime) {
		rateLimitRecord.Quota = 40
		rateLimitRecord.ResetTime = now.Add(1 * time.Minute)
	}

	if rateLimitRecord.Quota > 0 {
		rateLimitRecord.Quota -= quota
		return true
	}
	return false
}
