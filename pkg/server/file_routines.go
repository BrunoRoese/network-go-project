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
			buffer := make([]byte, 10000)
			n, _, err := newConn.ReadFromUDPAddrPort(buffer)

			go func() {
				if err != nil {
					slog.Error("[File saving] Error reading from UDP connection", slog.String("error", err.Error()))
					return
				}

				if n == 0 {
					slog.Warn("[File saving] Received empty message, skipping")
					return
				}

				req, err := parser.ParseLargeRequest(buffer[:n])
				if err != nil {
					slog.Error("[File saving] Error handling request", slog.String("error", err.Error()))
					return
				}

				if req.Information.Method == "END" {
					slog.Info("[File saving] Received end of file request")
					return
				}

				if err = validator.ValidateFileReq(req, currentChunk); err != nil {
					slog.Error("[File saving] Error validating request in file routine", slog.String("error", err.Error()))
					return
				}

				if inOrder, err := validator.CheckOrder(*req, currentChunk); !inOrder {
					if err != nil {
						slog.Error("[File saving] Error checking order", slog.String("error", err.Error()))
						req.Headers.XHeader["X-Chunk"] = strconv.Itoa(currentChunk)
						requests <- req
						return
					}
					slog.Info("[File saving] Ignoring chunk that is order")
				}

				if req.Information.Method == "CHUNK" {
					currentChunk++
					slog.Info("[File saving] Received chunk", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
				}
			}()

		}
	}()
}
