package seiya

import (
	"bytes"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type game struct {
	Name string
	URL  string
}

var (
	seiyaData   []game
	seiyaDataMu sync.RWMutex
)

func AllDataQuery() {

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

	logrus.Debugf("%+v", seiyaData)
	return nil
}
