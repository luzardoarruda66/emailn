package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	Db *gorm.DB
}

func (c *CampaignRepository) Create(campaign *campaign.Campaign) error {
	tx := c.Db.Create(campaign)
	return tx.Error
}

func (c *CampaignRepository) Update(campaign *campaign.Campaign) error {
	tx := c.Db.Save(campaign)
	return tx.Error
}

func (c *CampaignRepository) Get() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign
	tx := c.Db.Find(&campaigns)
	return campaigns, tx.Error
}

func (c *CampaignRepository) GetBy(id string) (*campaign.Campaign, error) {
	var campaigns campaign.Campaign
	tx := c.Db.Preload("Contacts").First(&campaigns, "id = ?", id)

	return &campaigns, tx.Error
}

func (c *CampaignRepository) Delete(campaign *campaign.Campaign) error {
	tx := c.Db.Delete(campaign)
	return tx.Error
}

func (c *CampaignRepository) GetCampaignsToBeSent() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign
	tx := c.Db.Preload("Contacts").Find(&campaigns, "status = ? and date_part('minute', now()::timestamp - updated_on::timestamp) > ?", campaign.Started, 1)
	return campaigns, tx.Error
}
