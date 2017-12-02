package downloader

import (
	"github.com/maxkulish/pageScan/config"
)

func CheckPageResponseChunk(chunk []config.Page) {

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
}

func checkPageResponse(page config.Page, ch chan config.Page, chFinished chan bool) error {

	defer func() {
		chFinished <- true
	}()

	page.RespCode = GetPageResponse(page.URL)

	ch <- page

	return error(nil)
}
