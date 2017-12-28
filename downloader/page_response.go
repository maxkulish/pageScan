package downloader

import (
	"time"

	"github.com/maxkulish/pageScan/config"
)

func CheckPageResponseChunk(chunk []config.Page) []config.Page {

	chPages := make(chan config.Page)
	chFinished := make(chan bool)

	results := make([]config.Page, 0, len(chunk))

	for _, page := range chunk {
		go checkPageResponse(page, chPages, chFinished)
	}

	for c := 0; c < len(chunk); {
		select {
		case page := <-chPages:
			results = append(results, page)
		case <-chFinished:
			c++
		}
	}

	return results
}

func checkPageResponse(page config.Page, ch chan config.Page, chFinished chan bool) error {

	defer func() {
		chFinished <- true
	}()

	start := time.Now()
	page.RespCode = GetPageResponse(page.URL)

	duration := time.Since(start)
	page.LoadTime = duration.Seconds()

	ch <- page

	return error(nil)
}
