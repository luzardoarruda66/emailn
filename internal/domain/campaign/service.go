package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/internalErrors"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(newCampaing contract.NewCampaign) (string, error) {

	campaign, err := NewCampaign(newCampaing.Name, newCampaing.Content, newCampaing.Emails)
	if err != nil {
		return "", err
	}

	err = s.Repository.Save(campaign)
	if err != nil {
		return "", internalerrors.ErrInternal
	}

	return campaign.ID, nil
}
