package database

import (
	"fmt"
	"log"
	"time"

	"github.com/maxkulish/pageScan/config"
)

func RetrieveSitemapPages(sitemap_id int) map[int]string {

	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(
		`SELECT id, page_url
			FROM sitemap_pages
			WHERE sitemap_link_id=$1`, sitemap_id)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Println("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func RetrieveUncheckedPages() map[int]string {
	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(`SELECT id, page_url
		FROM sitemap_pages
		WHERE load_time is NULL OR http_response = 0;`)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	pages := make(map[int]string)

	for rows.Next() {
		var pageID int
		var URL string

		if err := rows.Scan(&pageID, &URL); err != nil {
			log.Println("Can't read data from DB. Error: ", err)
		}

		pages[pageID] = URL
	}

	return pages
}

func BulkSavePagesResponse(pages []config.Page) bool {

	conn := Connection()
	defer conn.Close()

	rowQuery := ""

	dateTime := time.Now().UTC().Format(time.RFC3339)

	for _, page := range pages {
		rowQuery += fmt.Sprintf(
			"UPDATE sitemap_pages SET http_response=%d, "+
				"load_time=%f, updated=%s WHERE id = %d;",
			page.RespCode, dateTime, page.LoadTime, page.ID)
	}

	conn.Exec(rowQuery)

	return true
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
