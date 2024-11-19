package endpoints

import (
	"context"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignStart_200(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == "1"
	})).Return(nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("Patch", "/", nil)
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))

	rr := httptest.NewRecorder()

	_, status, err := handler.CampaignStart(rr, req)

	assert.Nil(err)
	assert.Equal(200, status)
}

func Test_CampaignStart_Err(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errExpected := errors.New("Something wrong")
	service.On("Start", mock.Anything).Return(errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("Patch", "/", nil)
	rr := httptest.NewRecorder()

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(err, errExpected)
}
