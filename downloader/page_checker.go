package downloader

import (
	"fmt"
	"log"
	"regexp"
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

func ExtractTitle(html string) []string {
	var titlePattern = regexp.MustCompile(`.*?<title[\sitemprop="name"\s]*?>(.*?)</title>`)

	title := titlePattern.FindStringSubmatch(html)

	return title

}

func ExtractHeaderOne(html string) []string {
	var headerOnePatt = regexp.MustCompile(`.*?<h1.*?>(.+?)</h1>`)

	header := headerOnePatt.FindStringSubmatch(html)

	return header
}

func ExtractDescription(html string) []string {
	var descrPattOne = regexp.MustCompile(`.*?<meta.*?name="description".*?content="(.+?)">`)

	descr := descrPattOne.FindStringSubmatch(html)

	if len(descr) == 0 {
		var descrPattTwo = regexp.MustCompile(`.*?<meta.*?content="(.+?)".*?name="description">`)
		descr = descrPattTwo.FindStringSubmatch(html)
	}

	return descr
}

func DownloadPageContent(pageURL string) {

	pageHTML, err := GetPageHtml(pageURL)

	if err != nil {
		log.Printf("[!] I can't find Title on page: %s\n", pageURL)
	}

	title := ExtractTitle(pageHTML)

	if len(title) == 0 {
		log.Println("<title> not found")
		fmt.Println(title)
	} else if len(title) > 2 {
		log.Printf("Found several <title> tags: %v", title)
		fmt.Println(title)
	} else {
		fmt.Println(title[len(title)-1])
	}

	header := ExtractHeaderOne(pageHTML)

	if len(header) == 0 {
		log.Println("<h1> not found")
		fmt.Println(header)
	} else if len(header) > 2 {
		log.Printf("Found several <h1> tags: %v", header)
		fmt.Println(header)
	} else {
		fmt.Println(header[len(header)-1])
	}

	description := ExtractDescription(pageHTML)

	if len(description) == 0 {
		log.Println("<meta name='description'> not found")
		fmt.Println(description)
	} else if len(description) > 2 {
		log.Printf("Found several <meta name='description'> tags: %v", header)
		fmt.Println(description)
	} else {
		fmt.Println(description[len(description)-1])
	}

}
