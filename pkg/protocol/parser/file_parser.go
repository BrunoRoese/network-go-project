package parser

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func ParseFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const chunkSize = 1024
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

func EncodeSha(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errors.New("Error opening file")
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", errors.New("Error hashing file")
	}

	hash := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hash)

	return hashHex, nil
}
