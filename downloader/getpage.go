package downloader

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/maxkulish/pageScan/agent"
)

func GetPageHtml(url string) (string, error) {
	timeout := time.Duration(20 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//log.Println("[!] Failed to crawl ", url)
		return "error", err
	}

	req.Header.Set("User-Agent", agent.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		//log.Println("[!] Failed to crawl ", url)
		return "error", err
	}
	defer resp.Body.Close()

	httpText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Println("[!] Failed to read server response ", url)
		return "error", err
	}

	return string(httpText), error(nil)
}

func GetPageResponse(url string) int {

	timeout := time.Duration(20 * time.Second)
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
