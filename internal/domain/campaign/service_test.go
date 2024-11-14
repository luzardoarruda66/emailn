package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internalErrors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Content min",
		Emails:    []string{"test@test.com"},
		CreatedBy: "teste@teste.com.br",
	}
	service = campaign.ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_CreateCampaign(t *testing.T) {
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repositoryMock

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_ValidadeRepositoryCreate(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to Create on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetBy(campaign.ID)

	assert.Equal(campaignReturned.ID, campaign.ID)
	assert.Equal(campaignReturned.Name, campaign.Name)
	assert.Equal(campaignReturned.Content, campaign.Content)
	assert.Equal(campaignReturned.Status, campaign.Status)
}
func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repositoryMock

	_, err := service.GetBy(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	campaignIdInvalid := "Invalid"
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock

	err := service.Delete(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal("campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problem(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body 1", []string{"teste@teste.com.br"}, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(errors.New("error to delete campaign"))
	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_has_success(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body 1", []string{"teste@teste.com.br"}, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(nil)
	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Nil(err)
}

func Test_Start_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	campaignIdInvalid := "Invalid"
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock

	err := service.Start(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock

	err := service.Start(campaign.ID)

	assert.Equal("campaign status invalid", err.Error())
}

func Test_Start_should_send_mail(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	emailWasSent := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignSaved.ID {
			emailWasSent = true
		}
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)
	assert.True(emailWasSent)
}

func Test_Start_ReturnError_when_func_SendMail_fail(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	service.Repository = repositoryMock
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignSaved.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaingToUpdate *campaign.Campaign) bool {
		return campaignSaved.ID == campaingToUpdate.ID && campaingToUpdate.Status == campaign.Done
	})).Return(nil)
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)
	assert.Equal(campaignSaved.Status, campaign.Done)
}
