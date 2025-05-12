package server

import (
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/validator"
	"log/slog"
	"net"
	"strconv"
)

func (s *Service) startFileSavingRoutine(newConn *net.UDPConn) {
	var currentChunk = 0
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

			if err = validator.ValidateFileReq(req, currentChunk); err != nil {
				slog.Error("[File saving] Error validating request in file routine", slog.String("error", err.Error()))
				continue
			}

			if inOrder, err := validator.CheckOrder(*req, currentChunk); !inOrder {
				if err != nil {
					slog.Error("[File saving] Error checking order", slog.String("error", err.Error()),
						slog.Int("expected", currentChunk+1),
						slog.String("received", req.Headers.XHeader["X-Chunk"]))
					req.Headers.XHeader["X-Chunk"] = strconv.Itoa(currentChunk)
					requests <- req
					continue
				}
				slog.Info("[File saving] Ignoring chunk that is out of order",
					slog.Int("expected", currentChunk+1),
					slog.String("received", req.Headers.XHeader["X-Chunk"]))
				continue
			}

			if req.Information.Method == "CHUNK" {
				currentChunk++
				req.Headers.XHeader["X-Chunk"] = strconv.Itoa(currentChunk)
				requests <- req
				slog.Info("[File saving] Received chunk", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
			}
		}
	}()
}
