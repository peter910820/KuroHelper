package bootstrap

import (
	"time"

	"github.com/sirupsen/logrus"

	"kurobidder/letao"
)

func KurobidderJob(baseURL string, filter letao.Filter, interval time.Duration, stopChan <-chan struct{}, dataChan chan<- []letao.AuctionItem) {
	logrus.Print("KurobidderJob 正在啟動...")
	ticker := time.NewTicker(interval * time.Minute)
	defer ticker.Stop()

	// 立即執行一次(只抓一筆)
	data, err := executeKurobidderCrawler(baseURL, filter)
	if err == nil && len(data) > 0 {
		data = data[:1]
		select {
		case dataChan <- data:
		case <-stopChan:
			return
		}
	}

	for {
		select {
		case <-ticker.C:
			data, err := executeKurobidderCrawler(baseURL, filter)
			if err == nil && len(data) > 0 {
				select {
				case dataChan <- data:
				case <-stopChan:
					return
				}
			}

		case <-stopChan:
			logrus.Println("KurobidderJob 正在關閉...")
			return
		}
	}
}

func executeKurobidderCrawler(baseURL string, filter letao.Filter) ([]letao.AuctionItem, error) {
	logrus.Info("開始執行kurobidder...")
	items, err := letao.LetaoCrawler(baseURL, filter)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.Info("kurobidder執行完成")
	return items, nil
}
