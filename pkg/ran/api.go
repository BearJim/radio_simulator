package ran

import (
	"context"
	"errors"

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
			RmState:          ue.RegisterState,
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
	return nil, errors.New("Not implemented")

}

func (a *apiService) Deregister(ctx context.Context, req *api.DeregisterRequest) (*api.DeregisterResponse, error) {
	return nil, errors.New("Not implemented")
}
