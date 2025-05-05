package protocol

type Protocol interface {
	Name() string
	BuildRequest(headers map[string]string, body string) string
}
