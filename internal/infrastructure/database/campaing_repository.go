package database

import "emailn/internal/domain/campaign"

type CampaingRepository struct {
	campaings []campaign.Campaign
}

func (c *CampaingRepository) Save(campaing *campaign.Campaign) error {
	c.campaings = append(c.campaings, *campaing)
	return nil
}
