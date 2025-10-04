package seiya

import (
	"bytes"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type game struct {
	Name string
	URL  string
}

var (
	seiyaData   []game
	seiyaDataMu sync.RWMutex
)

func GetGuideURL(keyword string) string {
	tokens := strings.Fields(strings.ToLower(keyword))
	weight := make(map[string]int)

	seiyaDataMu.RLock()
	defer seiyaDataMu.RUnlock()

	for _, seiya := range seiyaData {
		nameLower := strings.ToLower(seiya.Name)
		score := 0
		for _, token := range tokens {
			if strings.Contains(nameLower, token) {
				score++
			}
		}
		if score > 0 {
			weight[seiya.URL] = score
		}
	}

	// 選出最大權重
	var targetURL string
	maxValue := -1
	for k, v := range weight {
		if v > maxValue {
			targetURL = k
			maxValue = v
		}
	}

	if targetURL != "" {
		targetURL = "https://seiya-saiga.com/game/" + targetURL
	}
	return targetURL
}

func Init() error {
	r, err := sendGetRequest()
	if err != nil {
		return err
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return err
	}

	doc.Find(".table_hover").Each(func(i int, s *goquery.Selection) {
		s.Find("tbody tr td b a").Each(func(j int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			if exists {
				seiyaData = append(seiyaData, game{Name: a.Text(), URL: href})
			}
		})
	})

	return nil
}
