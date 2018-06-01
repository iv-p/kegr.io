package sync

import (
	"bytes"
	"context"
	"fmt"
	"log"

	grpc "google.golang.org/grpc"
	pbModel "kegr.io/protobuf/model/storage/server"
	pbServer "kegr.io/protobuf/server/storage"

	"kegr.io/storage_controller/model/keg"
	"kegr.io/storage_controller/model/liquid"
	"kegr.io/storage_controller/model/state"
)

type clientStatus int

const (
	ok       clientStatus = 0
	mismatch clientStatus = 1
	forward  clientStatus = 2
	down     clientStatus = 3
)

// InternalClient holds the information for another kegr instance
// in the cluster
type InternalClient struct {
	id      string
	address string
	client  pbServer.InternalClient
	conn    *grpc.ClientConn
	status  clientStatus
}

// IInternalClient is the InternalClient interface
type IInternalClient interface {
	GetServerInfo() *pbModel.ServerInfo
	GetID() string
	SetID(id string)
	GetAddress() string
	Shutdown()

	Ping(state state.IState) bool
	Register(ourID, ourAddress, ourPort string) ([]*pbModel.ServerInfo, error)
	GetState() map[string]keg.IKeg
	GetLiquid(kegID, liquidID string) liquid.ILiquid
}

// NewInternalClient initialises connection to the remote cerberus instance
func NewInternalClient(address string) IInternalClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pbServer.NewInternalClient(conn)

	return &InternalClient{
		address: address,
		client:  client,
		conn:    conn,
		status:  ok,
	}
}

func (c *InternalClient) GetServerInfo() *pbModel.ServerInfo {
	return &pbModel.ServerInfo{
		ID:      c.id,
		Address: c.address,
	}
}

func (c *InternalClient) Shutdown() {
	c.conn.Close()
}

// Ping pings the remote host to establish if they're still active
func (c *InternalClient) Ping(state state.IState) bool {
	s, err := state.GetHash()
	if err != nil {
		c.status = mismatch
		return false
	}

	prev := c.status
	if res, err := c.client.Ping(context.Background(), &pbServer.PingRequest{}); err == nil {
		if bytes.Equal(s, res.GetState()) {
			c.status = ok
		} else {
			c.status = mismatch
		}
	} else {
		c.status = down
	}

	if c.status == mismatch {
		log.Printf("state mismatch %v at %v\n", c.id, c.address)
	}

	if prev == ok && c.status == down {
		log.Printf("lost connection with %v at %v\n", c.id, c.address)
	}

	if prev == down && c.status == ok {
		log.Printf("reconnected with %v at %v\n", c.id, c.address)
	}

	if c.status == mismatch {
		return false
	}

	return true
}

// Register lets other instances on the cluster know we've joined
func (c *InternalClient) Register(ourID, ourAddress, ourPort string) ([]*pbModel.ServerInfo, error) {
	res, err := c.client.Register(context.Background(), &pbServer.RegisterRequest{
		Id:      ourID,
		Address: fmt.Sprintf("%v:%v", ourAddress, ourPort),
	})
	if err != nil {
		return nil, err
	}

	var others []*pbModel.ServerInfo
	for _, si := range res.Others {
		others = append(others, &pbModel.ServerInfo{ID: si.ID, Address: si.Address})
	}
	c.id = res.Id
	return others, nil
}

func (c *InternalClient) GetState() map[string]keg.IKeg {
	log.Printf("forcing recheck with %v at %v\n", c.id, c.address)
	state, err := c.client.GetState(context.Background(), &pbServer.GetStateRequest{})
	if err != nil {
		log.Println(err)
	}

	kegs := make(map[string]keg.IKeg)

	for id, k := range state.State.Kegs {
		kegs[id] = keg.FromProto(k)
	}

	return kegs
}

func (c *InternalClient) GetLiquid(kegID, liquidID string) liquid.ILiquid {
	log.Printf("getting file %v from keg %v from %v", liquidID, kegID, c.address)
	res, err := c.client.GetLiquid(
		context.Background(),
		&pbServer.GetLiquidRequest{
			KegId:    kegID,
			LiquidId: liquidID,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	return liquid.FromProto(res.Liquid)
}

func (c *InternalClient) GetID() string {
	return c.id
}

func (c *InternalClient) SetID(id string) {
	c.id = id
}

func (c *InternalClient) GetAddress() string {
	return c.address
}
