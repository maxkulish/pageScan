package database

import (
	"fmt"
	"time"

	"github.com/maxkulish/pageScan/config"
	log "gopkg.in/inconshreveable/log15.v2"
)

func RetrieveSitemapPages(sitemap_id int) map[int]string {

	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(
		`SELECT id, page_url
			FROM sitemap_pages
			WHERE sitemap_link_id=$1`, sitemap_id)
	if err != nil {
		log.Crit("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Info("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func RetrieveAllPages() map[int]string {
	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(`
		SELECT id, page_url
			FROM sitemap_pages;`)
	if err != nil {
		log.Crit("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Info("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func RetrieveUncheckedPages() map[int]string {
	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(`
		SELECT id, page_url
			FROM sitemap_pages
			WHERE load_time = 0.0 OR http_response = 0;`)
	if err != nil {
		log.Crit("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Info("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func RetrievePagesWithoutContent() map[int]string {

	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(`
		SELECT id, page_url
			FROM sitemap_pages
			WHERE
			title IS NULL OR h1 IS NULL
			OR description IS NULL OR title = 'not found';`)
	if err != nil {
		log.Crit("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Crit("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func BulkSavePagesResponse(pages []config.Page) bool {

	conn := Connection()
	defer conn.Close()

	rowQuery := ""
	var idList = make([]int, 0, len(pages))

	for _, page := range pages {

		query := fmt.Sprintf(
			"UPDATE sitemap_pages SET http_response=%d, "+
				"load_time=%f WHERE id = %d;",
			page.RespCode, page.LoadTime, page.ID)

		rowQuery += query
		idList = append(idList, page.ID)
	}

	conn.Exec(rowQuery)
	UpdatePageTime(idList)

	return true
}

func BulkUpdatePagesContent(pages []config.Page) bool {

	poolConn := PoolConnection()
	defer poolConn.Close()

	var idList = make([]int, 0, len(pages))

	for _, page := range pages {

		if _, err := poolConn.Exec("updateContent", page.Title, page.H1, page.Description, page.ID); err != nil {
			log.Error("Can't update page: %s", page.URL)
		}

		idList = append(idList, page.ID)
	}

	UpdatePageTime(idList)

	return true
}

func UpdatePageTime(idList []int) {
	conn := Connection()
	defer conn.Close()

	for _, pageId := range idList {
		conn.Exec("UPDATE sitemap_pages SET updated=$1 WHERE id = $2", time.Now(), pageId)
	}

}

func PagesToStruct(pages map[int]string) []config.Page {

	result := []config.Page{}

	for pageID, url := range pages {
		result = append(result, config.Page{
			ID:  pageID,
			URL: url,
		})
	}

	return result
}
