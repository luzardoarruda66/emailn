package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/gorm"
)

type CampaingRepository struct {
	Db *gorm.DB
}

func (c *CampaingRepository) Create(campaing *campaign.Campaign) error {
	tx := c.Db.Create(campaing)
	return tx.Error
}

func (c *CampaingRepository) Update(campaing *campaign.Campaign) error {
	tx := c.Db.Save(campaing)
	return tx.Error
}

func (c *CampaingRepository) Get() ([]campaign.Campaign, error) {
	var campaings []campaign.Campaign
	tx := c.Db.Find(&campaings)
	return campaings, tx.Error
}

func (c *CampaingRepository) GetBy(id string) (*campaign.Campaign, error) {
	var campaings campaign.Campaign
	tx := c.Db.Preload("Contacts").First(&campaings, "id = ?", id)

	return &campaings, tx.Error
}

func (c *CampaingRepository) Cancel(id string) error {
	var campaings campaign.Campaign
	tx := c.Db.First(&campaings, "id = ?", id)

	return tx.Error
}

func (c *CampaingRepository) Delete(campaign *campaign.Campaign) error {
	tx := c.Db.Delete(campaign)
	return tx.Error
}
