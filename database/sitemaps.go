package database

import "log"

func RetrieveAllSitemapsLinks() map[int]string {

	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query("SELECT id, url FROM sitemap_links")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	sitemaps := make(map[int]string)

	for rows.Next() {
		var sitemapID int
		var URL string

		if err := rows.Scan(&sitemapID, &URL); err != nil {
			log.Println("Can't read data from DB. Error: ", err)
		}

		sitemaps[sitemapID] = URL
	}

	return sitemaps
}

func RetrieveSitemapLinksByID(sitemapId int) map[int]string {

	conn := Connection()
	defer conn.Close()

	rows, err := conn.Query(`
			SELECT id, url
			FROM sitemap_links
			WHERE sitemap_id=$1
		`, sitemapId)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	sitemaps := make(map[int]string)

	for rows.Next() {
		var sitemapID int
		var URL string

		if err := rows.Scan(&sitemapID, &URL); err != nil {
			log.Println("Can't read data from DB. Error: ", err)
		}

		sitemaps[sitemapID] = URL
	}

	return sitemaps
}
