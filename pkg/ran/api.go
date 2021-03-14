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

func (a *apiService) GetUEs(ctx context.Context, params *api.GetUEsParams) (*api.GetUEsResponse, error) {
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

func (a *apiService) DescribeUE(ctx context.Context, params *api.DescribeUEParams) (*api.DescribeUEResponse, error) {
	return nil, errors.New("Not implemented")

}

func (a *apiService) Register(ctx context.Context, params *api.RegisterParams) (*api.RegisterResponse, error) {
	return nil, errors.New("Not implemented")

}

func (a *apiService) Deregister(ctx context.Context, params *api.DeregisterParams) (*api.DeregisterResponse, error) {
	return nil, errors.New("Not implemented")
}
