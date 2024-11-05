package internalmock

import (
	"emailn/internal/contract"

	"github.com/stretchr/testify/mock"
)

type CampaingServiceMock struct {
	mock.Mock
}

func (r *CampaingServiceMock) Create(newCampaing contract.NewCampaign) (string, error) {
	args := r.Called(newCampaing)
	return args.String(0), args.Error(1)
}

func (r *CampaingServiceMock) GetBy(id string) (*contract.CampaingResponse, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contract.CampaingResponse), args.Error(1)
}

func (r *CampaingServiceMock) Delete(id string) error {
	return nil
}
