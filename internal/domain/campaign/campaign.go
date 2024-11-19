package campaign

import (
	internalerrors "emailn/internal/internalErrors"
	"time"

	"github.com/rs/xid"
)

const (
	Pending = "Pending"
	Deleted = "Deleted"
	Started = "Started"
	Fail    = "Fail"
	Done    = "Done"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignId string
}

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    string    `gorm:"size:20"`
	CreatedBy string    `validate:"email" gorm:"size:50"`
}

func (c *Campaign) Done() {
	c.Status = Done
}

func (c *Campaign) Delete() {
	c.Status = Deleted
}

func (c *Campaign) Fail() {
	c.Status = Fail
}

func (c *Campaign) Started() {
	c.Status = Started
}

func NewCampaign(name string, content string, emails []string, createdBy string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))
	for i, value := range emails {
		contacts[i].Email = value
		contacts[i].ID = xid.New().String()
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
		Status:    Pending,
		CreatedBy: createdBy,
	}
	err := internalerrors.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}
	return nil, err

}
