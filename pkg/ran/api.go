package ran

import (
	"context"
	"errors"

	"github.com/free5gc/nas/security"
	"github.com/jay16213/radio_simulator/pkg/api"
)

type apiService struct {
	ranApp *RanApp
	api.UnimplementedAPIServiceServer
}

func (a *apiService) DescribeRAN(ctx context.Context, req *api.DescribeRANRequest) (*api.DescribeRANResponse, error) {
	resp := &api.DescribeRANResponse{
		Name: a.ranApp.Context().Name,
	}
	return resp, nil
}

func (a *apiService) GetUEs(ctx context.Context, req *api.GetUEsRequest) (*api.GetUEsResponse, error) {
	resp := &api.GetUEsResponse{}
	for _, ue := range a.ranApp.Context().UePool {
		resp.UeContexts = append(resp.UeContexts, &api.UEContext{
			Supi: ue.Supi,
			// Guti: ue.Guti,
			// CmState: ue.,
			RmState:          ue.RmState,
			NasUplinkCount:   ue.ULCount.Get(),
			NasDownlinkCount: ue.DLCount.Get(),
		})
	}
	return resp, nil
}

func (a *apiService) DescribeUE(ctx context.Context, req *api.DescribeUERequest) (*api.DescribeUEResponse, error) {
	return nil, errors.New("Not implemented")

}

func (a *apiService) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	ue := a.ranApp.Context().NewUE(req.Supi)
	ue.Supi = req.Supi
	ue.ServingPlmnId = req.ServingPlmn
	ue.AuthData.AuthMethod = req.AuthMethod
	ue.AuthData.K = req.K
	ue.AuthData.Opc = req.Opc
	ue.AuthData.Op = req.Op
	ue.AuthData.AMF = req.Amf
	ue.AuthData.SQN = req.Sqn
	switch req.CipheringAlg {
	case "NEA0":
		ue.CipheringAlg = security.AlgCiphering128NEA0
	case "NEA1":
		ue.CipheringAlg = security.AlgCiphering128NEA1
	case "NEA2":
		ue.CipheringAlg = security.AlgCiphering128NEA2
	case "NEA3":
		ue.CipheringAlg = security.AlgCiphering128NEA3
	}
	switch req.IntegrityAlg {
	case "NIA0":
		ue.IntegrityAlg = security.AlgIntegrity128NIA0
	case "NIA1":
		ue.IntegrityAlg = security.AlgIntegrity128NIA1
	case "NIA2":
		ue.IntegrityAlg = security.AlgIntegrity128NIA2
	case "NIA3":
		ue.IntegrityAlg = security.AlgIntegrity128NIA3
	}
	// amf selection
	ue.AMFEndpoint = a.ranApp.primaryAMFEndpoint
	a.ranApp.ngController.SendInitailUeMessage_RegistraionRequest(ue.AMFEndpoint, ue)
	// TODO: read register result
	return &api.RegisterResponse{Result: api.StatusCode_OK}, nil
}

func (a *apiService) Deregister(ctx context.Context, req *api.DeregisterRequest) (*api.DeregisterResponse, error) {
	return nil, errors.New("Not implemented")
}

func (a *apiService) SubscribeLog(req *api.LogStreamingRequest, stream api.APIService_SubscribeLogServer) error {
	supi := req.GetSupi()
	ue := a.ranApp.Context().FindUEBySupi(supi)
	if ue == nil {
		return errors.New("UE not found")
	}

	resp := &api.LogStreamingResponse{}
	for {
		resp.LogMessage = <-ue.ApiNotifyChan
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
