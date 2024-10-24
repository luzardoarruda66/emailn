package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infrastructure/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	campaingService := campaign.Service{
		Repository: &database.CampaingRepository{},
	}
	handler := endpoints.Handler{CampaingService: campaingService}
	r.Post("/campaigns", handler.CampaingPost)

	http.ListenAndServe(":3000", r)
}
