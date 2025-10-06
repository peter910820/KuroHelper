package seiya

import (
	"bytes"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type (
	game struct {
		Name string
		URL  string
	}

	candidate struct {
		Name  string
		URL   string
		Left  int
		Right int
	}
)

var rightWeightMap = map[string]struct{}{
	"full":    {},
	"voice":   {},
	"ver":     {},
	"the":     {},
	"edition": {},
	"hd":      {},
	"remake":  {},
	"plus":    {},
	"of":      {},
	"or":      {},
}

var (
	seiyaData   []game
	seiyaDataMu sync.RWMutex
)

func GetGuideURL(keyword string) string {
	tokens := strings.Fields(strings.ToLower(keyword))
	var candidateGames []candidate

	seiyaDataMu.RLock()
	defer seiyaDataMu.RUnlock()

	for _, seiya := range seiyaData {
		nameLower := strings.ToLower(seiya.Name)
		leftWeight := 0
		rightWeight := 0
		for _, token := range tokens {
			token = strings.Map(func(r rune) rune {
				switch r {
				case '-', '!', '！', '.', '～':
					return -1 // delete rune
				}
				return r
			}, token)
			if strings.Contains(nameLower, token) {
				_, ok := rightWeightMap[token]
				if ok {
					rightWeight++
				} else {
					leftWeight++
				}
			}
		}
		if leftWeight > 0 {
			candidateGames = append(candidateGames, candidate{
				Name:  seiya.Name,
				URL:   seiya.URL,
				Left:  leftWeight,
				Right: rightWeight,
			})
		}
	}

	var targetURL string
	if len(candidateGames) == 0 {
		return targetURL
	}

	// 選出最大權重
	sort.Slice(candidateGames, func(i, j int) bool {
		if candidateGames[i].Left != candidateGames[j].Left {
			return candidateGames[i].Left > candidateGames[j].Left
		}
		return candidateGames[i].Right > candidateGames[j].Right
	})

	if !strings.HasPrefix(candidateGames[0].URL, "https://") && strings.TrimSpace(candidateGames[0].URL) != "" {
		targetURL = "https://seiya-saiga.com/game/" + candidateGames[0].URL
	} else {
		targetURL = candidateGames[0].URL
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
