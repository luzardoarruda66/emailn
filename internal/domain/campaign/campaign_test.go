package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name     = "Campaign X"
	content  = "Content min"
	contacts = []string{"email1@e.com", "email2@e.com"}
	fake     = faker.New()
)

func Test_NewCampaign_CreateCampaing(t *testing.T) {
	assert := assert.New(t)

	campaing, _ := NewCampaign(name, content, contacts)

	assert.Equal(campaing.Name, name)
	assert.Equal(campaing.Content, content)
	assert.Equal(len(campaing.Contacts), len(contacts))
}

func Test_NewCampaign_IDIsNotNill(t *testing.T) {
	assert := assert.New(t)

	campaing, _ := NewCampaign(name, content, contacts)

	assert.NotNil(campaing.ID)
}

func Test_NewCampaign_CreatedOnIsMustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaing, _ := NewCampaign(name, content, contacts)

	assert.Greater(campaing.CreatedOn, now)
}

func Test_NewCampaing_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts)
	assert.Equal("name is required with min 5", err.Error())
}

func Test_NewCampaing_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts)

	assert.Equal("name is required with max 24", err.Error())
}

func Test_NewCampaing_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "", contacts)
	assert.Equal("content is required with min 5", err.Error())
}

func Test_NewCampaing_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts)
	assert.Equal("content is required with max 1024", err.Error())
}

func Test_NewCampaing_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{})
	assert.Equal("contacts is required with min 1", err.Error())
}

func Test_NewCampaing_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"email_invalid"})
	assert.Equal("email is invalid", err.Error())
}
