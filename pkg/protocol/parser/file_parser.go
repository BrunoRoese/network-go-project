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
	return encodeFileToBase64Chunks(filePath, 1024)
}

func encodeFileToBase64Chunks(filePath string, chunkSize int) ([]string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(file)

	var chunks []string
	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunks = append(chunks, encoded[i:end])
	}

	return chunks, nil
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
