package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/maxkulish/pageScan/config"
	"github.com/maxkulish/pageScan/database"
	"github.com/maxkulish/pageScan/downloader"
	"github.com/maxkulish/pageScan/utils"
)

func main() {

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Example: ./%s -resp -speed=15\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	var speed int
	flag.IntVar(
		&speed,
		"speed",
		15,
		"Scan speed. Example: ./pageScanLinux -speed=15")

	pageResp := flag.Bool("resp", false, "Check pages response")
	sitemapId := flag.Int("sitemap", 0, "Check by sitemap ID")

	flag.Parse()

	globalStart := time.Now()
	log.Printf("\033[92m[+] Global start: %s\033[0m\n", globalStart)

	// export environment variables from .env
	config.SetEnvironment()

	if *sitemapId != 0 {
		checkResponseForSitemap(speed, *sitemapId)
	}

	if *pageResp == true {
		checkResponseUncheckedPages(speed)
	}

	log.Printf("\033[92m[+] Done! Spent %s\033[0m\n", time.Since(globalStart))

}

func checkResponseForSitemap(chunkSize, sitemapId int) {

	sitemaps := database.RetrieveSitemapLinksByID(sitemapId)
	log.Printf("Get %d sitemaps to check Response Code\n", len(sitemaps))

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

	chunks := utils.ChunkifyPages(pagesToCheck, chunkSize)

	for _, chunk := range chunks {
		results := downloader.CheckPageResponseChunk(chunk)

		for _, page := range results {
			log.Printf("Code: %d. LoadTime: %f URL: %s\n", page.RespCode, page.LoadTime, page.URL)
		}

		database.BulkSavePagesResponse(results)
	}

}

func checkResponseUncheckedPages(chunkSize int) {

	pages := database.RetrieveUncheckedPages()

	if len(pages) == 0 {
		log.Println("There are no pages to check")
		os.Exit(1)
	}

	pagesToCheck := database.PagesToStruct(pages)

	chunks := utils.ChunkifyPages(pagesToCheck, chunkSize)

	for _, chunk := range chunks {
		results := downloader.CheckPageResponseChunk(chunk)

		for _, page := range results {
			log.Printf("Code: %d. LoadTime: %f\tURL: %s\n", page.RespCode, page.LoadTime, page.URL)
		}

		database.BulkSavePagesResponse(results)
	}

}
