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
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service = campaign.ServiceImp{
		Repository: repositoryMock,
	}
	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func Test_Create_Campaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_CreateCampaign(t *testing.T) {
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setUp()
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_ValidadeRepositoryCreate(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to Create on database"))

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(campaignReturned.ID, campaignPending.ID)
	assert.Equal(campaignReturned.Name, campaignPending.Name)
	assert.Equal(campaignReturned.Content, campaignPending.Content)
	assert.Equal(campaignReturned.Status, campaignPending.Status)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy("invalid_campaiign")

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal("campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problem(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_has_success(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(err)
}

func Test_Start_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)

	assert.Equal("campaign status invalid", err.Error())
}

func Test_Start_should_send_mail(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	emailWasSent := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPending.ID {
			emailWasSent = true
		}
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPending.ID)
	assert.True(emailWasSent)
}

func Test_Start_ReturnError_when_func_SendMail_fail(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaingToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaingToUpdate.ID && campaingToUpdate.Status == campaign.Done
	})).Return(nil)
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPending.ID)
	assert.Equal(campaignPending.Status, campaign.Done)
}
