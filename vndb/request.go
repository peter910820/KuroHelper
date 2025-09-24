package vndb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func sendRequest(apiRoute string, jsonBytes []byte) ([]byte, error) {
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
