package vndb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type rateLimitStruct struct {
	Quota     int
	ResetTime time.Time
}

var (
	rateLimitRecord = rateLimitStruct{
		Quota:     40,
		ResetTime: time.Now().Add(1 * time.Minute),
	}
	rateLimitMu sync.RWMutex
)

func sendRequest(apiRoute string, jsonBytes []byte) ([]byte, error) {
	if !rateLimit(1) {
		return nil, fmt.Errorf("Quota exhausted")
	}
	resp, err := http.Post(os.Getenv("VNDB_ENDPOINT")+apiRoute, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func rateLimit(quota int) bool {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

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
