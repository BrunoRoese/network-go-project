package client

type Client struct {
	Ip            string
	Port          int
	LastHeartbeat int64
}
