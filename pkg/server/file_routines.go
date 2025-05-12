package server

import (
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/validator"
	"log/slog"
	"net"
	"strconv"
)

var (
	lastRecChunk = -1
)

func (s *Service) startFileSavingRoutine(newConn *net.UDPConn) {
	go func() {
		for {
			slog.Info("[File saving] File saving routine started on source", newConn.LocalAddr().(*net.UDPAddr).String())
			buffer := make([]byte, 1024)
			n, _, err := newConn.ReadFromUDPAddrPort(buffer)

			if err != nil {
				slog.Error("[File saving] Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				slog.Error("[File saving] Error handling request", slog.String("error", err.Error()))
				continue
			}

			err = validator.ValidateFileReq(req, lastRecChunk)

			if err != nil {
				slog.Error("[File saving] Error validating request in file routine", slog.String("error", err.Error()))
				continue
			}

			if req.Information.Method == "CHUNK" {
				chunk, err := strconv.Atoi(req.Headers.XHeader["X-Chunk"])

				if err != nil {
					slog.Error("[File saving] Error converting chunk to int", slog.String("error", err.Error()))
					continue
				}

				lastRecChunk = chunk
				slog.Info("[File saving] Received chunk", slog.String("chunk", req.Headers.XHeader["X-Chunk"]))
			}

			req.Headers.XHeader["X-Chunk"] = strconv.Itoa(lastRecChunk)

			requests <- req
		}
	}()
}
