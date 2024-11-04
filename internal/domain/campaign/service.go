package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/internalErrors"
)

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaingResponse, error)
}
type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	if err != nil {
		return "", err
	}
	err = s.Repository.Save(campaign)
	if err != nil {
		return "", internalerrors.ErrInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaingResponse, error) {
	campaing, err := s.Repository.GetBy(id)
	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return &contract.CampaingResponse{
		ID:      campaing.ID,
		Name:    campaing.Name,
		Content: campaing.Content,
		Status:  campaing.Status}, nil
}
