package erogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	internalerrors "kurohelper/errors"
	erogsmodels "kurohelper/models/erogs"
)

func GetCreatorByFuzzy(search string) (*erogsmodels.FuzzySearchCreatorResponse, error) {
	formData := url.Values{}
	formData.Set("sql", buildFuzzySearchCreatorSQL(search))

	resp, err := http.Post(os.Getenv("EROGS_ENDPOINT"), "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", resp.StatusCode)
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	selection := doc.Find("td").First() // 只取第一個符合的
	jsonText := selection.Text()

	if strings.TrimSpace(jsonText) == "" {
		return nil, internalerrors.ErrVndbNoResult
	}

	var res erogsmodels.FuzzySearchCreatorResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	if len(res.Games) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	return &res, nil
}
