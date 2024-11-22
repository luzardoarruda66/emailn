package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/infrastructure/database"
	"emailn/internal/infrastructure/mail"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error to load .env file")
	}

	db := database.NewDB()
	repository := database.CampaignRepository{Db: db}
	campaignService := campaign.ServiceImp{
		Repository: &repository,
		SendMail:   mail.SendMail,
	}

	for {
		campaigns, err := repository.GetCampaignsToBeSent()

		println("Ämount of campaigns to be sent:", len(campaigns))

		if err != nil {
			println(err.Error())
		}

		for _, campaign := range campaigns {
			campaignService.SendEmailAndUpdateStatus(&campaign)
		}
		time.Sleep(10 * time.Second)
	}

}
