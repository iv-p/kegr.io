package sync

import (
	"fmt"
	"log"
	"time"

	pbModel "kegr.io/protobuf/model/storage/server"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/state"
)

// SyncService holds information of all the connected clients
// in the cluster
type SyncService struct {
	ISyncService

	id      string
	clients map[string]IInternalClient
	ss      state.IStateService
}

// ISyncService is the SyncService interface
type ISyncService interface {
	GetID() string
	Register(clusterMember string)
	GetClientAddresses() []*pbModel.ServerInfo
	AddClient(id, address string)
}

// NewSyncService takes a single cluster member and registers
// recursively with it and any other cluster member it
// gets.
func NewSyncService(ss state.IStateService) *SyncService {
	serv := &SyncService{
		id:      config.C.MachineName,
		ss:      ss,
		clients: make(map[string]IInternalClient),
	}

	go serv.monitor()

	if len(config.C.Other) > 0 {
		serv.Register(config.C.Other)
	}

	return serv
}

// GetID returns the current server's id
func (ss *SyncService) GetID() string {
	return ss.id
}

// Register initialises connection with a cluster
func (ss *SyncService) Register(clusterMember string) {
	var instances []*pbModel.ServerInfo

	log.Printf("trying to connect to cluster %v\n", clusterMember)

	client := NewInternalClient(clusterMember)
	instances, _ = client.Register(ss.id, config.C.Address, config.C.InternalGrpcPort)

	ss.addClientToMap(client.GetID(), client)

	for _, instance := range instances {
		if _, exist := ss.clients[instance.ID]; !exist {
			ss.Register(instance.Address)
		}
	}
}

// GetClientAddresses returns the server information for all clients
// currently connected
func (ss *SyncService) GetClientAddresses() []*pbModel.ServerInfo {
	var others []*pbModel.ServerInfo
	for _, v := range ss.clients {
		others = append(others, v.GetServerInfo())
	}

	return others
}

// AddClient adds
func (ss *SyncService) AddClient(id, address string) {
	if _, exist := ss.clients[id]; exist {
		return
	}
	client := NewInternalClient(address)
	client.SetID(id)
	ss.addClientToMap(id, client)
}

func (ss *SyncService) addClientToMap(id string, client IInternalClient) {
	log.Printf("connected to client %v at %v\n", client.GetID(), client.GetAddress())
	ss.clients[id] = client
}

func (ss *SyncService) monitor() {
	for range time.Tick(1 * time.Second) {
		for _, client := range ss.clients {
			if ok := client.Ping(ss.ss.GetState()); !ok {
				ss.forceRecheck(client)
			}
		}
	}
}

func (ss *SyncService) forceRecheck(client IInternalClient) {
	otherState := client.GetState()
	diff := ss.ss.Diff(otherState)

	for kegID, kegDiff := range diff {
		if kegDiff.Options != nil {
			ss.ss.UpdateKeg(kegID, kegDiff.Options)
		}

		for _, liquidInfo := range kegDiff.Content {
			liquid := client.GetLiquid(kegID, string(liquidInfo.GetID()))

			if keg, err := ss.ss.GetKegByID(kegID); err == nil {
				path := fmt.Sprintf("%s/%s", config.C.DataRoot, kegID)
				if err = liquid.ToFile(path); err == nil {
					keg.AddLiquid(liquid.GetLiquidInfo())
				}
			}
		}
	}
}
