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

	db := database.NewDB()
	campaingService := campaign.ServiceImp{
		Repository: &database.CampaingRepository{Db: db},
	}
	handler := endpoints.Handler{CampaingService: &campaingService}
	r.Post("/campaigns", endpoints.HandlerError(handler.CampaingPost))
	r.Get("/campaigns/{id}", endpoints.HandlerError(handler.CampaingGetById))
	r.Delete("/campaigns/delete/{id}", endpoints.HandlerError(handler.CampaingDelete))

	http.ListenAndServe(":3000", r)
}
