package main

import (
	"flag"
	"fmt"

	"os"
	"path/filepath"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/maxkulish/pageScan/config"
	"github.com/maxkulish/pageScan/database"
	"github.com/maxkulish/pageScan/downloader"
	"github.com/maxkulish/pageScan/utils"
)

func main() {

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fileName := filepath.Base(os.Args[0])
		fmt.Printf(
			"Examples\nCheck responses: ./%s -resp -speed=15\nDownload pages: ./%s -content\nRescan all pages: ./%s -content -rescan\n",
			fileName, fileName, fileName)
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
	content := flag.Bool("content", false, "Download page title, h1, description")
	rescan := flag.Bool("rescan", false, "Rescan response time or page title, h1, description")

	flag.Parse()

	globalStart := time.Now()

	// export environment variables from .env
	config.SetEnvironment()

	if *sitemapId != 0 {
		checkResponseForSitemap(speed, *sitemapId)
	}

	if *pageResp == true {

		if *rescan == true {
			checkPagesResponse(speed, true)
		} else {
			checkPagesResponse(speed, false)
		}

	}

	if *content == true {

		if *rescan == true {
			downloadPagesContent(speed, true)
		} else {
			downloadPagesContent(speed, false)
		}

	}

	log.Info("[+] Done! Spent ", "time", time.Since(globalStart))
}

func checkResponseForSitemap(chunkSize, sitemapId int) {

	sitemaps := database.RetrieveSitemapLinksByID(sitemapId)
	log.Info("Checking Response Code", "number", len(sitemaps))

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

		database.BulkSavePagesResponse(results)
	}

}

// Get pages from Database and get page response: 200, 301, 404, 500
func checkPagesResponse(chunkSize int, rescan bool) {

	var pages = make(map[int]string)

	if rescan {
		pages = database.RetrieveAllPages()
	} else {
		pages = database.RetrieveUncheckedPages()
	}

	log.Info("Checking Response Code", "pages", len(pages))

	if len(pages) == 0 {
		log.Warn("There are no pages to check")
		os.Exit(1)
	}

	pagesToCheck := database.PagesToStruct(pages)

	chunks := utils.ChunkifyPages(pagesToCheck, chunkSize)

	for _, chunk := range chunks {
		results := downloader.CheckPageResponseChunk(chunk)

		database.BulkSavePagesResponse(results)
	}
}

// Get pages from Database and download h1, title, description
func downloadPagesContent(chunkSize int, rescan bool) {

	var pages = make(map[int]string)
	if rescan {
		pages = database.RetrieveAllPages()
	} else {
		pages = database.RetrievePagesWithoutContent()
	}

	if len(pages) == 0 {
		log.Warn("There are no pages to check")
		os.Exit(1)
	}

	log.Info("Downloading Pages Content", "pages", len(pages))

	pagesToCheck := database.PagesToStruct(pages)

	chunks := utils.ChunkifyPages(pagesToCheck, chunkSize)

	for _, chunk := range chunks {

		results := downloader.DownloadPageContentChunk(chunk)

		database.BulkUpdatePagesContent(results)
	}

}
