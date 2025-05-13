package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/validator"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// FileWriter manages writing chunks to a file in a streaming manner
type FileWriter struct {
	file          *os.File
	mutex         sync.Mutex
	currentChunk  int
	requestId     string
	totalChunks   int
	resourcesPath string
	writtenChunks []string
}

var chunkMap = make(map[int]string)

// NewFileWriter creates a new FileWriter instance
func NewFileWriter(requestId string) (*FileWriter, error) {
	resourcesPath := "resources"
	if _, err := os.Stat(resourcesPath); os.IsNotExist(err) {
		if err := os.MkdirAll(resourcesPath, 0755); err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(resourcesPath, requestId+".pdf")
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	slog.Info("[File saving] Created file", slog.String("path", filePath))

	return &FileWriter{
		file:          file,
		mutex:         sync.Mutex{},
		currentChunk:  0,
		requestId:     requestId,
		totalChunks:   0,
		resourcesPath: resourcesPath,
		writtenChunks: []string{},
	}, nil
}

func (fw *FileWriter) WriteChunk(req *protocol.Request) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	for _, check := range fw.writtenChunks {
		if check == req.Headers.XHeader["X-Checksum"] {
			slog.Info("[File saving] Chunk already written, skipping", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
			return nil
		}
	}

	chunkStr := req.Headers.XHeader["X-Chunk"]
	chunk, err := strconv.Atoi(chunkStr)
	slog.Info("Chunk received to write: ", slog.String("n", chunkStr))
	if err != nil {
		return err
	}

	if endStr, ok := req.Headers.XHeader["X-End"]; ok {
		if end, err := strconv.Atoi(endStr); err == nil {
			fw.totalChunks = end
		}
	}

	if chunk != fw.currentChunk+1 && chunk != 0 {
		return nil
	}

	decodedData, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return err
	}

	_, err = fw.file.Write(decodedData)
	if err != nil {
		return err
	}

	fw.writtenChunks = append(fw.writtenChunks, req.Headers.XHeader["X-Checksum"])
	checksum := req.Headers.XHeader["X-Checksum"]
	calculated := sha256.Sum256([]byte(req.Body))
	calculatedStr := hex.EncodeToString(calculated[:])
	if checksum != calculatedStr {
		return fmt.Errorf("[File saving] Checksum mismatch on chunk %d", fw.currentChunk+1)
	}

	fw.currentChunk = chunk

	slog.Info("[File saving] Wrote chunk to file",
		slog.Int("chunk", chunk),
		slog.Int("totalChunks", fw.totalChunks),
		slog.String("requestId", fw.requestId))

	// If this is the last chunk, close the file
	if fw.totalChunks > 0 && fw.currentChunk >= fw.totalChunks {
		slog.Info("[File saving] All chunks received, closing file",
			slog.String("requestId", fw.requestId))
		return fw.file.Close()
	}

	return nil
}

// Close closes the file
func (fw *FileWriter) Close() error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	if fw.file != nil {
		return fw.file.Close()
	}
	return nil
}

func (s *Service) startFileSavingRoutine(newConn *net.UDPConn) {
	var fileWriter *FileWriter
	var fileWriterMutex sync.Mutex
	chunks := make([]int, 0)
	expectedChunk := 0

	go func() {
		for {
			slog.Info("[File saving] Waiting for message on port", slog.Int("port", newConn.LocalAddr().(*net.UDPAddr).Port))
			buffer := make([]byte, 2048)
			n, _, err := newConn.ReadFromUDPAddrPort(buffer)

			if err != nil {
				slog.Error("[File saving] Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			if n == 0 {
				slog.Warn("[File saving] Received empty message, skipping")
				continue
			}

			req, err := parser.ParseLargeRequest(buffer[:n])
			if err != nil {
				slog.Error("[File saving] Error handling request", slog.String("error", err.Error()))
				continue
			}

			currentChunk, err := strconv.Atoi(req.Headers.XHeader["X-Chunk"])

			var receivedChunk bool
			for _, c := range chunks {
				if c == currentChunk {
					slog.Info("[File saving] Chunk already received, skipping", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
					receivedChunk = true
				}
			}

			if receivedChunk {
				slog.Info("[File saving] Chunk already received, skipping", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
				requests <- req
				continue
			}

			if currentChunk != expectedChunk {
				slog.Error("[File saving] Chunk out of order", slog.Int("expected", expectedChunk), slog.Int("received", currentChunk))
				chunkMap[currentChunk] = req.Body
				chunks = append(chunks, currentChunk)
				continue
			}

			if err != nil {
				slog.Error("[File saving] Error converting chunk to int", slog.String("error", err.Error()))
			}

			if req.Information.Method == "END" {
				fileWriterMutex.Lock()
				if fileWriter != nil {
					if err := fileWriter.Close(); err != nil {
						slog.Error("[File saving] Error closing file", slog.String("error", err.Error()))
					}
					fileWriter = nil
				}
				fileWriterMutex.Unlock()

				resourcesPath := "resources"
				filePath := filepath.Join(resourcesPath, req.Information.Id.String()+".pdf")
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					slog.Error("[File saving] File does not exist", slog.String("path", filePath))
					return
				}

				encodedSha, err := parser.EncodeSha(filePath)

				if err != nil {
					slog.Error("[File saving] Error encoding SHA", slog.String("error", err.Error()))
					continue
				}

				if encodedSha != req.Body {
					slog.Error("[File saving] SHA mismatch", slog.String("expected", req.Body), slog.String("calculated", encodedSha))
					req.Information.Method = "NACK"
				}

				requests <- req
				slog.Info("[File saving] Received END request, file closed")
				return
			}

			if err = validator.ValidateFileReq(req, currentChunk); err != nil {
				slog.Error("[File saving] Error validating request in file routine", slog.String("error", err.Error()))
				continue
			}

			if req.Information.Method == "CHUNK" {
				fileWriterMutex.Lock()
				if fileWriter == nil {
					var initErr error
					fileWriter, initErr = NewFileWriter(req.Information.Id.String())
					if initErr != nil {
						slog.Error("[File saving] Error initializing file writer", slog.String("error", initErr.Error()))
						fileWriterMutex.Unlock()
						continue
					}
				}
				fileWriterMutex.Unlock()
				expectedChunk++
				chunks = append(chunks, currentChunk)

				if err := fileWriter.WriteChunk(req); err != nil {
					slog.Error("[File saving] Error writing chunk to file", slog.String("error", err.Error()))
				}

				requests <- req
				slog.Info("[File saving] Received chunk", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
			}
		}
	}()
}
