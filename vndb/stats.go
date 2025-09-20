package vndb

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetStats() ([]byte, error) {
	resp, err := http.Get(os.Getenv("VNDB_ENDPOINT") + "/stats")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", resp.StatusCode)
	}

	return body, nil
}
