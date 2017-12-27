package downloader

import (
	"log"
	"net/http"
	"time"

	"github.com/maxkulish/pageScan/agent"
	"golang.org/x/net/html"
)

func GetParsedHTML(url string) (*html.Node, error) {
	timeout := time.Duration(40 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[!] Failed to crawl ", url)
		return nil, err
	}

	req.Header.Set("User-Agent", agent.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[!] Failed to crawl ", url)
		return nil, err
	}

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	return root, error(nil)
}

func GetPageResponse(url string) int {

	timeout := time.Duration(40 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[!] Failed to crawl ", url)
		return 0
	}

	req.Header.Set("User-Agent", agent.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[!] Failed to crawl ", url)
		return 0
	}

	defer resp.Body.Close()

	return resp.StatusCode
}
