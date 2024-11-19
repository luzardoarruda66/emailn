package endpoints

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignStart_200(t *testing.T) {
	setup()
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == "1"
	})).Return(nil)
	req, rr := newHttpTest("PATCH", "/", nil)
	req = addParameter(req, "id", "1")

	_, status, err := handler.CampaignStart(rr, req)

	assert.Nil(t, err)
	assert.Equal(t, 200, status)
}

func Test_CampaignStart_Err(t *testing.T) {
	setup()
	errExpected := errors.New("Something wrong")
	service.On("Start", mock.Anything).Return(errExpected)
	req, rr := newHttpTest("PATCH", "/", nil)

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(t, err, errExpected)
}
