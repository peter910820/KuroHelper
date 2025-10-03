package seiya

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"

	kurohelpererrors "kurohelper/errors"
)

func sendGetRequest() (*goquery.Document, error) {
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

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
