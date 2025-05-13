package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/protocol/validator"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

type FileService struct {
	FilePath    string
	conn        *net.UDPConn
	clientAddr  *net.UDPAddr
	currentId   uuid.UUID
	encodedFile string

	currentChunk     chan int
	stopSending      chan bool
	receivedResponse []int
}

func (s *FileService) StartTransfer(ip string, filePath string) error {
	if err := validator.Validate(ip, filePath); err != nil {
		slog.Error("Error validating input", slog.String("error", err.Error()))
		return err
	}
	s.FilePath = filePath
	s.stopSending = make(chan bool)
	s.receivedResponse = make([]int, 0)

	specifiedClient := client.FindByIp(ip)

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return nil
	}

	conn, err := network.CreateConn()

	if err != nil {
		slog.Error("Error creating connection", slog.String("error", err.Error()))
		return err
	}
	defer conn.Close()

	s.conn = conn

	fileContent, err := parser.ParseFile(s.FilePath)
	encodedSha, err := parser.EncodeSha(s.FilePath)

	if err != nil {
		slog.Error("Error parsing file", slog.String("error", err.Error()))
		return err
	}

	s.encodedFile = encodedSha

	err = s.signalStart(specifiedClient)

	if err != nil {
		slog.Error("Error signaling start", slog.String("error", err.Error()))
		s.stopSending <- true
	}

	s.startRoutines(fileContent)

	<-s.stopSending

	s.close()

	return nil
}

func (s *FileService) signalStart(specifiedClient *client.Client) error {
	headers := map[string]string{
		"X-Chunk": "0",
	}
	req := (&protocol.File{}).BuildRequest(headers, s.FilePath, *s.conn.LocalAddr().(*net.UDPAddr))

	s.currentId = req.Information.Id

	jsonReq, err := json.Marshal(req)

	if err != nil {
		slog.Error("Error marshaling jsonReq")
		return err
	}

	_, _ = network.SendRequest(specifiedClient.Ip, specifiedClient.Port, jsonReq)

	return nil
}

func (s *FileService) startRoutines(fileContent []string) {
	s.currentChunk = make(chan int)

	s.startSendingRoutine(fileContent)
	s.startListeningRoutine(fileContent)
}

func (s *FileService) startSendingRoutine(fileContent []string) {
	go func(fileContent []string) {
		for chunk := range s.currentChunk {
			if chunk < len(fileContent)-1 {
				currentChunk := fileContent[chunk]
				//slog.Info("Sending chunk", "chunk", currentChunk)

				checksumBytes := sha256.Sum256([]byte(currentChunk))
				checksum := hex.EncodeToString(checksumBytes[:])

				headers := map[string]string{
					"X-Chunk":    strconv.Itoa(chunk),
					"X-Checksum": checksum,
					"X-End":      strconv.Itoa(len(fileContent) - 1),
					"requestId":  s.currentId.String(),
				}

				res, err := parser.ParseProtocol(&protocol.Chunk{}, s.conn, currentChunk, headers)

				if err != nil {
					//slog.Error("Error marshalling request", slog.String("error", err.Error()))
					continue
				}

				for retry := 0; retry < 20; retry++ {
					_, _ = network.SendRequest(s.clientAddr.IP.String(), s.clientAddr.Port, res)
					send := false

					for _, c := range s.receivedResponse {
						if c == chunk {
							send = true
						}
					}

					if send {
						//slog.Info("Chunk sent successfully", "chunk", chunk)
						break
					} else {
						//slog.Error("Chunk not sent, retrying", "chunk", chunk)
					}
					time.Sleep(200 * time.Millisecond)

					if retry == 19 {
						chunk++
						s.currentChunk <- chunk
					}
				}
			} else {
				//slog.Error("Chunk index out of bounds", "chunk", chunk)

				headers := map[string]string{
					"requestId": s.currentId.String(),
				}

				res, err := parser.ParseProtocol(&protocol.End{}, s.conn, s.encodedFile, headers)

				if err != nil {
					//slog.Error("Error marshalling request", slog.String("error", err.Error()))
					continue
				}

				for retry := 0; retry < 20; retry++ {
					_, _ = network.SendRequest(s.clientAddr.IP.String(), s.clientAddr.Port, res)

					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}(fileContent)
}

func (s *FileService) startListeningRoutine(fileContent []string) {
	//slog.Info("Starting listening routine")
	chunkTrack := 0
	go func(fileContent []string) {
		for {
			receivedCount := len(s.receivedResponse)
			targetEnd := len(fileContent) - 1

			if targetEnd > 0 {
				percentage := (float64(receivedCount) / float64(targetEnd)) * 100
				displayProgressBar(percentage)
			}
			buffer := make([]byte, 1024)
			n, _, err := s.conn.ReadFromUDPAddrPort(buffer)
			if err != nil {
				//slog.Error("Error reading from UDP", slog.String("error", err.Error()))
				return
			}

			//slog.Info("Response received", "data", string(buffer[:n]), "from", addr.String())
			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				//slog.Error("Error parsing request", slog.String("error", err.Error()))
				return
			}

			if req.Information.Id != s.currentId {
				//slog.Error("Invalid request ID", slog.String("expected", s.currentId.String()), slog.String("received", req.Information.Id.String()))
				return
			}

			ip, port, err := parser.ParseSource(req.Information.Source)

			if err != nil {
				//slog.Error("Error parsing source", slog.String("error", err.Error()))
				return
			}

			if req.Information.Method == "NACK" {
				//slog.Error("Received NACK or END, stopping sending", slog.String("method", req.Information.Method))
				displayProgressBar(100)
				time.Sleep(150 * time.Millisecond)
				slog.Info("Received NACK, file transfer failed")
				s.stopSending <- true
				return
			}

			if req.Information.Method == "END" {
				displayProgressBar(100)
				time.Sleep(150 * time.Millisecond)
				slog.Info("Received END, file transfer completed")
				s.stopSending <- true
				return
			}

			//slog.Info("Parsed source", "ip", ip, "port", port)

			if s.clientAddr == nil {
				s.clientAddr = &net.UDPAddr{
					IP:   net.ParseIP(ip),
					Port: port,
				}
			}

			//slog.Info("Client address set", "ip", s.clientAddr.IP.String(), "port", s.clientAddr.Port)
			go func() {
				receivedChunk, err := strconv.Atoi(req.Headers.XHeader["X-Chunk"])

				if err != nil {
					//slog.Error("Error converting chunk to int", slog.String("error", err.Error()))

					s.currentChunk <- 0

					return
				}

				if receivedChunk == chunkTrack {
					chunkTrack++
				} else {
					for _, chunk := range s.receivedResponse {
						if chunk == receivedChunk {
							//slog.Info("Chunk already received, skipping", "chunk", receivedChunk)
							return
						}
						if chunk > receivedChunk {
							//slog.Info("Chunk already received, skipping", "chunk", receivedChunk)
							return
						}
					}
					//slog.Info("Chunk out of order, resending", "chunk", receivedChunk, "expected", chunkTrack)
					chunkTrack++
				}

				//slog.Info("receivedChunk", "chunk", receivedChunk, "chunkTrack", chunkTrack)

				s.receivedResponse = append(s.receivedResponse, receivedChunk)

				s.currentChunk <- chunkTrack
			}()
		}
	}(fileContent)
}

func displayProgressBar(percentage float64) {
	barWidth := 50
	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}

	progress := int(percentage / 100 * float64(barWidth))
	bar := "[" + strings.Repeat("=", progress) + strings.Repeat(" ", barWidth-progress) + "]"
	fmt.Printf("\r%s %.2f%%", bar, percentage)
}

func (s *FileService) close() {
	close(s.currentChunk)
	close(s.stopSending)
}
