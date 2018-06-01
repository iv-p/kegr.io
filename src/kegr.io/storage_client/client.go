package storage_client

import (
	"log"

	grpc "google.golang.org/grpc"
	pbServer "kegr.io/protobuf/server/storage"
)

// SingleClient holds the information for any other cerberus instance
// in the cluster
type SingleClient struct {
	address string
	client  pbServer.ExternalClient
	conn    *grpc.ClientConn
}

// ISingleClient is the SingleClient interface
type ISingleClient interface {
	Shutdown()
	Get() pbServer.ExternalClient
}

// NewSingleClientClient initialises connection to the remote cerberus instance
func NewSingleClientClient(address string) *SingleClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("fail to dial: %v", err)
	}
	client := pbServer.NewExternalClient(conn)

	return &SingleClient{
		address: address,
		client:  client,
		conn:    conn,
	}
}

func (c *SingleClient) Get() pbServer.ExternalClient {
	return c.client
}

// Shutdown closes the connection with that server
func (c *SingleClient) Shutdown() {
	c.conn.Close()
}
