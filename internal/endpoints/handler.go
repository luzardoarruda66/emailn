package endpoints

import "emailn/internal/domain/campaign"

type Handler struct {
	CampaingService campaign.Service
}
