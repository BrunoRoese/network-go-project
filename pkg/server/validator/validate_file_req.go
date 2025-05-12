package validator

import (
	"errors"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/google/uuid"
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

	if req.Information.Method == "CHUNK" && lastChunk == reqChunk {
		return errors.New("Duplicated chunk found")
	}

	return nil
}
