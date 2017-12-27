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

	rows, err := conn.Query(`
		SELECT id, page_url
			FROM sitemap_pages
			WHERE load_time = 0.0 OR http_response = 0;`)
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

func BulkSavePagesContent(pages []config.Page) bool {

	conn := Connection()
	defer conn.Close()

	rowQuery := ""
	var idList = make([]int, 0, len(pages))

	for _, page := range pages {

		query := fmt.Sprintf(
			"UPDATE sitemap_pages SET title=%s, "+
				"h1=%s, description=%s WHERE id = %d;",
			page.Title, page.H1, page.Description, page.ID)

		rowQuery += query
		idList = append(idList, page.ID)
	}

	conn.Exec(rowQuery)
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
