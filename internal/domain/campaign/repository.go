package campaign

type Repository interface {
	Save(campaing *Campaign) error
	Get() ([]Campaign, error)
	GetBy(id string) (*Campaign, error)
}
