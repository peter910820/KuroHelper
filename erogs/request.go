package erogs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	kurohelpererrors "kurohelper/errors"
)

func sendPostRequest(sql string) (string, error) {
	formData := url.Values{}
	formData.Set("sql", sql)

	req, err := http.NewRequest("POST", os.Getenv("EROGS_ENDPOINT"), strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return "", err
	}

	selection := doc.Find("td").First() // 只取第一個符合的
	jsonText := selection.Text()

	if strings.TrimSpace(jsonText) == "" {
		return "", kurohelpererrors.ErrSearchNoContent
	}

	return jsonText, nil
}
