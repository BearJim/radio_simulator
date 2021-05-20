package ran

import (
	"context"
	"errors"
	"fmt"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/security"
	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_packet"
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
			NasUplinkCount:   ue.ULCount.ToUint32(),
			NasDownlinkCount: ue.DLCount.ToUint32(),
		})
	}
	return resp, nil
}

func (a *apiService) DescribeUE(ctx context.Context, req *api.DescribeUERequest) (*api.DescribeUEResponse, error) {
	ue := a.ranApp.Context().FindUEBySupi(req.GetSupi())
	if ue == nil {
		return nil, fmt.Errorf("UE not found (supi: %s)", req.GetSupi())
	}

	return &api.DescribeUEResponse{UeContext: &api.UEContext{
		Supi:             ue.Supi,
		RmState:          ue.RmState,
		CmState:          ue.CmState,
		AmfUeNgapId:      ue.AmfUeNgapId,
		RanUeNgapId:      ue.RanUeNgapId,
		NasUplinkCount:   ue.ULCount.ToUint32(),
		NasDownlinkCount: ue.DLCount.ToUint32(),
	}}, nil
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
	ue.FollowOnRequest = req.FollowOnRequest
	a.ranApp.ngController.NewNASConnection(ue)
	a.ranApp.ngController.SendInitailUeMessage_RegistraionRequest(ue.AMFEndpoint, ue)

	// wait result
	result := <-ue.ApiNotifyChan
	return &api.RegisterResponse{
		StatusCode: result.Status,
		Body:       result.Message,
		UeContext: &api.UEContext{
			Supi:             ue.Supi,
			RmState:          ue.RmState,
			CmState:          ue.CmState,
			NasUplinkCount:   ue.ULCount.ToUint32(),
			NasDownlinkCount: ue.DLCount.ToUint32(),
			AmfUeNgapId:      ue.AmfUeNgapId,
			RanUeNgapId:      ue.RanUeNgapId,
		},
		RestartCount:     int32(result.RestartCount),
		RestartTimestamp: result.RestartTimeStamp.UnixNano(),
	}, nil
}

func (a *apiService) Deregister(ctx context.Context, req *api.DeregisterRequest) (*api.DeregisterResponse, error) {
	ue := a.ranApp.ctx.FindUEBySupi(req.GetSupi())
	if ue == nil {
		return nil, fmt.Errorf("UE not found (supi: %s)", req.GetSupi())
	}

	nasPdu, err := nas_packet.GetDeregistrationRequest(ue, 0) //normoal release
	if err != nil {
		logger.ApiLog.Error(err.Error())
		return &api.DeregisterResponse{StatusCode: api.StatusCode_ERROR}, nil
	}
	a.ranApp.ngController.SendUplinkNASTransport(ue.AMFEndpoint, ue, nasPdu)

	// wait result
	result := <-ue.ApiNotifyChan
	return &api.DeregisterResponse{StatusCode: result.Status, Body: result.Message}, nil
}

func (a *apiService) ServiceRequestProc(ctx context.Context, req *api.ServiceRequest) (*api.ServiceRequestResult, error) {
	ue := a.ranApp.ctx.FindUEBySupi(req.GetSupi())
	if ue == nil {
		return nil, fmt.Errorf("UE not found (supi: %s)", req.GetSupi())
	}

	serviceType := uint8(0)
	switch req.ServiceType {
	case api.ServiceType_Signalling:
		serviceType = nasMessage.ServiceTypeSignalling
	case api.ServiceType_Data:
		serviceType = nasMessage.ServiceTypeData
	default:
		serviceType = nasMessage.ServiceTypeSignalling
	}

	logger.ApiLog.Infow("Service Request Procedure", "supi", ue.Supi, "id", ue.AmfUeNgapId, "rid", ue.RanUeNgapId)
	a.ranApp.ngController.NewNASConnection(ue)
	nasPdu, err := nas_packet.GetServiceRequest(ue, serviceType)
	if err != nil {
		logger.ApiLog.Error(err.Error())
		return nil, fmt.Errorf("build error: %+v", err)
	}
	a.ranApp.ngController.SendInitailUeMessage(ue.AMFEndpoint, ue, ue.GutiStr[7:], nasPdu)

	// wait result
	result := <-ue.ApiNotifyChan
	return &api.ServiceRequestResult{StatusCode: result.Status, Body: result.Message}, nil
}

func (a *apiService) SubscribeLog(req *api.LogStreamingRequest, stream api.APIService_SubscribeLogServer) error {
	supi := req.GetSupi()
	ue := a.ranApp.Context().FindUEBySupi(supi)
	if ue == nil {
		return errors.New("UE not found")
	}

	// resp := &api.LogStreamingResponse{}
	// for {
	// 	resp.LogMessage = <-ue.ApiNotifyChan
	// 	if err := stream.Send(""); err != nil {
	// 		logger.ApiLog.Errorf("err: %+v", err)
	// 		return err
	// 	}
	// }
	return errors.New("Not implemented")
}
