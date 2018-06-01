package server

import (
	"context"
	"fmt"
	"time"

	pbKeg "kegr.io/protobuf/model/storage/keg"
	pbLiquid "kegr.io/protobuf/model/storage/liquid"
	pbServer "kegr.io/protobuf/server/storage"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/model/keg"
	"kegr.io/storage_controller/model/liquid"
	"kegr.io/storage_controller/state"
	"kegr.io/storage_controller/util"
)

// Server defines the grpc service
type ExternalServer struct {
	ss state.IStateService
}

func NewExternalServer(ss state.IStateService) *ExternalServer {
	return &ExternalServer{
		ss: ss,
	}
}

// CreateLiquid returns the merkle tree of this server
func (es *ExternalServer) CreateLiquid(ctx context.Context, req *pbServer.CreateLiquidRequest) (*pbServer.CreateLiquidResponse, error) {
	liquid := liquid.FromProto(req.GetLiquid())
	liquid.SetID(util.ID())
	liquid.SetLastUpdated(time.Now().Unix())

	keg, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.CreateLiquidResponse{}, err
	}

	err = liquid.ToFile(fmt.Sprintf("%s/%s", config.C.DataRoot, keg.GetID()))
	if err != nil {
		return &pbServer.CreateLiquidResponse{}, err
	}

	err = keg.AddLiquid(liquid.GetLiquidInfo())
	return &pbServer.CreateLiquidResponse{}, err
}

// GetLiquid returns the merkle tree of this server
func (es *ExternalServer) GetLiquid(ctx context.Context, req *pbServer.GetLiquidRequest) (*pbServer.GetLiquidResponse, error) {
	_, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.GetLiquidResponse{}, err
	}

	liquid, err := liquid.FromFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))

	if err != nil {
		return &pbServer.GetLiquidResponse{}, err
	}

	return &pbServer.GetLiquidResponse{
		Liquid: liquid.ToProto(),
	}, nil
}

// UpdateLiquid updates a liquid
func (es *ExternalServer) UpdateLiquid(ctx context.Context, req *pbServer.UpdateLiquidRequest) (*pbServer.UpdateLiquidResponse, error) {
	l := req.GetLiquid()
	liquid := liquid.FromProto(l)
	liquid.SetLastUpdated(time.Now().Unix())

	keg, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.UpdateLiquidResponse{}, err
	}

	err = liquid.ToFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))
	if err != nil {
		return &pbServer.UpdateLiquidResponse{}, err
	}

	err = keg.UpdateLiquid(liquid.GetLiquidInfo())
	return &pbServer.UpdateLiquidResponse{}, err
}

// UpdateLiquidOptions updates the options of a liquid
func (es *ExternalServer) UpdateLiquidOptions(ctx context.Context, req *pbServer.UpdateLiquidOptionsRequest) (*pbServer.UpdateLiquidOptionsResponse, error) {
	options := liquid.OptionsFromProto(req.GetOptions())
	keg, err := es.ss.GetKegByID(req.GetKegId())

	if err != nil {
		return &pbServer.UpdateLiquidOptionsResponse{}, err
	}

	liquid, err := liquid.FromFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))
	if err != nil {
		return &pbServer.UpdateLiquidOptionsResponse{}, err
	}

	liquid.SetLastUpdated(time.Now().Unix())
	liquid.SetOptions(options)

	err = liquid.ToFile(fmt.Sprintf("%s/%s", config.C.DataRoot, req.GetKegId()))
	if err != nil {
		return &pbServer.UpdateLiquidOptionsResponse{}, err
	}

	err = keg.UpdateLiquid(liquid.GetLiquidInfo())
	return &pbServer.UpdateLiquidOptionsResponse{}, err
}

// DeleteLiquid marks the liquid as deleted
func (es *ExternalServer) DeleteLiquid(ctx context.Context, req *pbServer.DeleteLiquidRequest) (*pbServer.DeleteLiquidResponse, error) {
	keg, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.DeleteLiquidResponse{}, err
	}

	liquid, err := liquid.FromFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))
	if err != nil || liquid.IsDeleted() {
		return &pbServer.DeleteLiquidResponse{}, err
	}

	liquid.SetLastUpdated(time.Now().Unix())
	liquid.SetDeleted(true)

	err = liquid.ToFile(fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, req.GetKegId(), req.GetLiquidId(), config.C.LiquidExtension))
	if err != nil {
		return &pbServer.DeleteLiquidResponse{}, err
	}

	err = keg.UpdateLiquid(liquid.GetLiquidInfo())
	return &pbServer.DeleteLiquidResponse{}, err
}

// CreateKeg returns the merkle tree of this server
func (es *ExternalServer) CreateKeg(ctx context.Context, req *pbServer.CreateKegRequest) (*pbServer.CreateKegResponse, error) {
	options := keg.OptionsFromProto(req.GetOptions())
	keg, err := es.ss.CreateKeg(options)
	if err != nil {
		return &pbServer.CreateKegResponse{}, err
	}
	return &pbServer.CreateKegResponse{
		KegID: keg.GetID(),
	}, nil
}

// GetKeg returns the merkle tree of this server
func (es *ExternalServer) GetKeg(ctx context.Context, req *pbServer.GetKegRequest) (*pbServer.GetKegResponse, error) {
	keg, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.GetKegResponse{}, err
	}
	return &pbServer.GetKegResponse{
		Options: keg.GetOptions().ToProto(),
	}, nil
}

// GetKegs returns all kegs on server
func (es *ExternalServer) GetKegs(ctx context.Context, req *pbServer.GetKegsRequest) (*pbServer.GetKegsResponse, error) {
	kegs := es.ss.GetKegs()

	protoKegs := make(map[string]*pbKeg.Info)

	for id, keg := range kegs {
		protoKegs[id] = keg.GetInfo().ToProto()
	}

	return &pbServer.GetKegsResponse{
		Kegs: protoKegs,
	}, nil
}

// GetKegLiquids returns all liquids in that keg
func (es *ExternalServer) GetKegLiquids(ctx context.Context, req *pbServer.GetKegLiquidsRequest) (*pbServer.GetKegLiquidsResponse, error) {
	keg, err := es.ss.GetKegByID(req.GetKegId())
	if err != nil {
		return &pbServer.GetKegLiquidsResponse{}, err
	}

	liquids := keg.GetLiquids()
	var infos []*pbLiquid.Info

	for _, liq := range liquids {
		infos = append(infos, liq.ToProto())
	}

	return &pbServer.GetKegLiquidsResponse{
		Liquids: infos,
	}, nil
}

// UpdateKegOptions updates the options of a Keg
func (es *ExternalServer) UpdateKegOptions(ctx context.Context, req *pbServer.UpdateKegOptionsRequest) (*pbServer.UpdateKegOptionsResponse, error) {
	options := keg.OptionsFromProto(req.GetOptions())
	es.ss.UpdateKeg(req.GetKegId(), options)
	return &pbServer.UpdateKegOptionsResponse{}, nil
}

// DeleteKeg marks the Keg as deleted
func (es *ExternalServer) DeleteKeg(ctx context.Context, req *pbServer.DeleteKegRequest) (*pbServer.DeleteKegResponse, error) {
	err := es.ss.DeleteKeg(req.GetKegId())
	return &pbServer.DeleteKegResponse{}, err
}
