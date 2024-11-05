package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/internalErrors"
	"errors"
)

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaingResponse, error)
	Cancel(id string) error
	Delete(id string) error
}
type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	if err != nil {
		return "", err
	}
	err = s.Repository.Create(campaign)
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
		ID:                   campaing.ID,
		Name:                 campaing.Name,
		Content:              campaing.Content,
		Status:               campaing.Status,
		AmountOfEmailsToSend: len(campaing.Contacts),
	}, nil
}

func (s *ServiceImp) Cancel(id string) error {
	campaing, err := s.Repository.GetBy(id)
	if err != nil {
		return internalerrors.ErrInternal
	}

	if campaing.Status != Pending {
		return errors.New("campaign status invalid")
	}
	campaing.Cancel()
	err = s.Repository.Update(campaing)
	if err != nil {
		return internalerrors.ErrInternal
	}
	return nil
}

func (s *ServiceImp) Delete(id string) error {

	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerrors.ErrInternal
	}

	if campaign.Status != Pending {
		return errors.New("campaign status invalid")
	}

	campaign.Delete()
	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}
	return nil
}
