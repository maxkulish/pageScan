package downloader

import (
	"strings"

	"log"
	"time"

	"github.com/maxkulish/pageScan/config"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExtractTitle(htmlBody *html.Node) string {

	// Search for the title
	title, ok := scrape.Find(htmlBody, scrape.ByTag(atom.Title))
	if ok {
		return strings.TrimSpace(scrape.Text(title))
	}

	return "not found"

}

func ExtractHeaderOne(htmlBody *html.Node) string {

	h1, ok := scrape.Find(htmlBody, scrape.ByTag(atom.H1))

	if ok {
		return strings.TrimSpace(scrape.Text(h1))
	}

	return "not found"
}

func ExtractDescription(htmlBody *html.Node) string {

	descr := scrape.FindAll(htmlBody, scrape.ByTag(atom.Meta))

	for _, meta := range descr {
		res := scrape.Attr(meta, "name") == "description"

		if res {
			found := scrape.Attr(meta, "content")
			return strings.TrimSpace(found)
		}
	}

	return "not found"
}

func DownloadPageContentChunk(chunk []config.Page) []config.Page {

	chPages := make(chan config.Page)
	chFinished := make(chan bool)

	results := make([]config.Page, 0, len(chunk))

	for _, page := range chunk {
		go DownloadPageContent(page, chPages, chFinished)
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

func DownloadPageContent(page config.Page, ch chan config.Page, chFinished chan bool) error {

	defer func() {
		chFinished <- true
	}()

	start := time.Now()
	htmlBody, err := GetParsedHTML(page.URL)
	if err != nil {
		log.Printf("I can't donwload page HTML. URL: %s", page.URL)
		chFinished <- true
	}
	duration := time.Since(start)
	page.LoadTime = duration.Seconds()

	page.Title = ExtractTitle(htmlBody)
	page.H1 = ExtractHeaderOne(htmlBody)
	page.Description = ExtractDescription(htmlBody)

	ch <- page

	return error(nil)

}
