package database

import "log"

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
