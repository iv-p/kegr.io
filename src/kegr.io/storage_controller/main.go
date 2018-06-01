package main

import (
	"fmt"
	"log"
	"net"

	"kegr.io/protobuf/server/storage"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/server"
	"kegr.io/storage_controller/state"
	"kegr.io/storage_controller/sync"

	"google.golang.org/grpc"
)

func main() {
	config.Load()

	stateService := state.NewStateService()
	syncService := sync.NewSyncService(stateService)

	grpcInternalServer := grpc.NewServer()
	internalServer := server.NewInternalServer(syncService, stateService)
	storage.RegisterInternalServer(grpcInternalServer, internalServer)

	log.Println("starting grpc servers")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.C.InternalGrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go grpcInternalServer.Serve(lis)

	grpcExternalServer := grpc.NewServer()
	externalServer := server.NewExternalServer(stateService)
	storage.RegisterExternalServer(grpcExternalServer, externalServer)

	lis, err = net.Listen("tcp", fmt.Sprintf(":%v", config.C.ExternalGrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcExternalServer.Serve(lis)
}
