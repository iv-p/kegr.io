package storage_client

import pbServer "kegr.io/protobuf/server/storage"

// Client handles the connection to a storage controller cluster
type Client struct {
	clients map[string]ISingleClient
}

// IClient is the Client interface
type IClient interface {
	Get() pbServer.ExternalClient
}

// NewClient initialises a new Client object and connects to a storage
// controller cluster
func NewClient(address string) *Client {
	cli := &Client{
		clients: make(map[string]ISingleClient),
	}
	cli.clients["temp"] = NewSingleClientClient(address)
	return cli
}

// Get returns the best storage controller client to use
func (c *Client) Get() pbServer.ExternalClient {
	return c.clients["temp"].Get()
}
