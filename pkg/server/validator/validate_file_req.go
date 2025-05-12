package validator

import (
	"errors"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

func ValidateFileReq(req *protocol.Request, lastChunk int) error {
	if req.Information.Id == uuid.Nil {
		return errors.New("request ID is empty")
	}

	if req.Body == "" {
		return errors.New("request body is empty")
	}

	if req.Information.Method != "CHUNK" && req.Information.Method != "END" {
		return errors.New("Invalid method")
	}

	reqChunk, err := strconv.Atoi(req.Headers.XHeader["X-Chunk"])

	if err != nil {
		return errors.New("Invalid chunk number")
	}

	slog.Info("[File saving] Validating request", slog.String("chunk", req.Headers.XHeader["X-Chunk"]), slog.Int("lastChunk", lastChunk))

	if lastChunk == 0 && reqChunk == 0 {
		return nil
	}

	if req.Information.Method == "CHUNK" && lastChunk == reqChunk {
		return errors.New("Duplicated chunk found")
	}

	return nil
}

func CheckOrder(req protocol.Request, lastChunk int) (bool, error) {
	reqChunk, _ := strconv.Atoi(req.Headers.XHeader["X-Chunk"])

	if reqChunk == lastChunk+1 {
		return true, nil
	}

	if reqChunk < lastChunk {
		slog.Info("[File saving] Ignoring chunk that is too far behind", slog.Int("expected", lastChunk+1), slog.Int("received", reqChunk))
		return false, nil
	}

	if reqChunk > lastChunk+1 {
		return false, errors.New("Chunk too far ahead, expected " + strconv.Itoa(lastChunk+1) + " but got " + strconv.Itoa(reqChunk))
	}

	return true, nil
}
