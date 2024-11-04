package campaign

type Repository interface {
	Create(campaing *Campaign) error
	Update(campaing *Campaign) error
	Get() ([]Campaign, error)
	GetBy(id string) (*Campaign, error)
	Delete(campaing *Campaign) error
}
