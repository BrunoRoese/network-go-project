package parser

import (
	"encoding/base64"
	"io"
	"os"
)

func ParseFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const chunkSize = 512
	buffer := make([]byte, chunkSize)
	var base64Chunks []string

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			base64Chunk := base64.StdEncoding.EncodeToString(buffer[:n])
			base64Chunks = append(base64Chunks, base64Chunk)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return base64Chunks, nil
}
