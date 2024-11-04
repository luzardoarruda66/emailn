package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/internalErrors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Create(campaing *Campaign) error {
	args := r.Called(campaing)
	return args.Error(0)
}

func (r *repositoryMock) Update(campaing *Campaign) error {
	args := r.Called(campaing)
	return args.Error(0)
}

func (r *repositoryMock) Get() ([]Campaign, error) {
	return nil, nil
}

func (r *repositoryMock) GetBy(id string) (*Campaign, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Campaign), nil
}

func (r *repositoryMock) Delete(campaing *Campaign) error {
	args := r.Called(campaing)
	return args.Error(0)
}

var (
	newCampaign = contract.NewCampaign{
		Name:    "Test Y",
		Content: "Content min",
		Emails:  []string{"test@test.com"},
	}
	service = ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repository := new(repositoryMock)
	repository.On("Save", mock.Anything).Return(nil)
	service.Repository = repository

	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_SaveCampaign(t *testing.T) {
	repository := new(repositoryMock)
	repository.On("Save", mock.MatchedBy(func(campaing *Campaign) bool {
		if campaing.Name != newCampaign.Name ||
			campaing.Content != newCampaign.Content ||
			len(campaing.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repository

	service.Create(newCampaign)

	repository.AssertExpectations(t)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_ValidadeRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repository := new(repositoryMock)
	repository.On("Save", mock.Anything).Return(errors.New("error to save on database"))
	service.Repository = repository

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaing(t *testing.T) {
	assert := assert.New(t)
	campaing, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repository := new(repositoryMock)
	repository.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaing.ID
	})).Return(campaing, nil)
	service.Repository = repository

	campaingReturned, _ := service.GetBy(campaing.ID)

	assert.Equal(campaingReturned.ID, campaing.ID)
	assert.Equal(campaingReturned.Name, campaing.Name)
	assert.Equal(campaingReturned.Content, campaing.Content)
	assert.Equal(campaingReturned.Status, campaing.Status)
}
func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaing, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repository := new(repositoryMock)
	repository.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repository

	_, err := service.GetBy(campaing.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}
