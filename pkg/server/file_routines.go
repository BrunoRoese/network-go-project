package server

import (
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/validator"
	"log/slog"
	"net"
	"strconv"
)

var (
	lastRecChunk = map[string]int{}
)

func (s *Service) startFileSavingRoutine(newConn *net.UDPConn) {
	go func() {
		for {
			slog.Info("[File saving] Waiting for message on port", slog.Int("port", newConn.LocalAddr().(*net.UDPAddr).Port))
			buffer := make([]byte, 10000)
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

			var currentChunk int
			if req.Information.Method == "CHUNK" {
				chunk, err := strconv.Atoi(req.Headers.XHeader["X-Chunk"])
				if err != nil {
					slog.Error("[File saving] Error converting chunk to int", slog.String("error", err.Error()))
					continue
				}
				currentChunk = chunk
			}

			if err = validator.ValidateFileReq(req, lastRecChunk[req.Information.Id.String()]); err != nil {
				slog.Error("[File saving] Error validating request in file routine", slog.String("error", err.Error()))
				continue
			}

			if inOrder, err := validator.CheckOrder(*req, lastRecChunk[req.Information.Id.String()]); !inOrder {
				if err != nil {
					slog.Error("[File saving] Error checking order", slog.String("error", err.Error()))
					req.Headers.XHeader["X-Chunk"] = strconv.Itoa(lastRecChunk[req.Information.Id.String()])
					requests <- req
					continue
				}
				slog.Info("[File saving] Ignoring chunk that is order")
			}

			if req.Information.Method == "CHUNK" {
				lastRecChunk[req.Information.Id.String()] = currentChunk
				slog.Info("[File saving] Received chunk", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
			}
		}
	}()
}
