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

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
