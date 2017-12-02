package main

import (
	"fmt"
	"log"
	"time"

	"github.com/maxkulish/pageScan/config"
	"github.com/maxkulish/pageScan/database"
	"github.com/maxkulish/pageScan/utils"
)

func main() {

	globalStart := time.Now()
	log.Printf("\033[92m[+] Global start: %s\033[0m", globalStart)

	// export environment variables from .env
	config.SetEnvironment()

	checkResponseAllPages()

	log.Printf("\033[92m[+] Done! Spent %s\033[0m", time.Since(globalStart))

}

func checkResponseAllPages() {

	sitemaps := database.RetrieveSitemapLinksByID(4)
	fmt.Println(sitemaps)

	pagesToCheck := []config.Page{}
	for stmpId, _ := range sitemaps {

		pages := database.RetrieveSitemapPages(stmpId)

		for pageID, url := range pages {
			pagesToCheck = append(pagesToCheck, config.Page{
				ID:  pageID,
				URL: url,
			})
		}
	}

	res := utils.ChunkifyPages(pagesToCheck, 10)

	fmt.Println(len(res))
	//fmt.Printf("Code: %d. Page: %s\n", downloader.GetPageResponse(url), url)

}
