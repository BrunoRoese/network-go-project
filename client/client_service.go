package client

type Client struct {
	ip   string
	port int
}

func CreateClient(ip string, port int) *Client {
	return &Client{ip: ip, port: port}
}
