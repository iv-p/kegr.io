package server

import (
	"context"
	"fmt"

	pb "kegr.io/protobuf/server/storage"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/model/liquid"
	"kegr.io/storage_controller/state"
	"kegr.io/storage_controller/sync"
)

// InternalServer defines the grpc service
type InternalServer struct {
	is sync.ISyncService
	ss state.IStateService
}

// NewInternalServer returns an initialised internal server object
func NewInternalServer(is sync.ISyncService, ss state.IStateService) *InternalServer {
	return &InternalServer{
		is: is,
		ss: ss,
	}
}

// Ping is a simple helathcheck endpoint
func (is *InternalServer) Ping(ctx context.Context, ping *pb.PingRequest) (*pb.PingResponse, error) {
	hash, err := is.ss.GetHash()
	return &pb.PingResponse{
		State: hash,
	}, err
}

// Register is the healthcheck other instances use to make sure
// they're on this instances radar
func (is *InternalServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	others := is.is.GetClientAddresses()
	is.is.AddClient(req.Id, req.Address)
	return &pb.RegisterResponse{
		Response: "ok",
		Id:       is.is.GetID(),
		Others:   others,
	}, nil
}

// GetState returns the merkle tree of this server
func (is *InternalServer) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	return &pb.GetStateResponse{
		State: is.ss.GetState().ToProto(),
	}, nil
}

// GetPeers returns the merkle tree of this server
func (is *InternalServer) GetPeers(ctx context.Context, req *pb.GetPeersRequest) (*pb.GetPeersResponse, error) {
	return &pb.GetPeersResponse{
		Peers: is.is.GetClientAddresses(),
	}, nil
}

// GetLiquid returns a liquid
func (is *InternalServer) GetLiquid(ctx context.Context, req *pb.GetLiquidRequest) (*pb.GetLiquidResponse, error) {
	_, err := is.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pb.GetLiquidResponse{}, err
	}

	liquid, err := liquid.FromFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))
	if err != nil {
		return &pb.GetLiquidResponse{}, err
	}

	return &pb.GetLiquidResponse{
		Liquid: liquid.ToProto(),
	}, err
}
