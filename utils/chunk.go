package utils

import "github.com/maxkulish/pageScan/config"

func ChunkifyPages(array []config.Page, chunkSize int) [][]config.Page {

	chunk := make([]config.Page, 0, chunkSize)
	chunks := make([][]config.Page, 0, len(array)/chunkSize+1)

	for len(array) >= chunkSize {
		chunk, array = array[:chunkSize], array[chunkSize:]
		chunks = append(chunks, chunk)
	}

	if len(array) > 0 {
		chunks = append(chunks, array[:len(array)])
	}

	return chunks
}
