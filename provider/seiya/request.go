package seiya

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	kurohelpererrors "kurohelper/errors"
)

func sendGetRequest() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("SEIYA_ENDPOINT"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return r, nil
}
