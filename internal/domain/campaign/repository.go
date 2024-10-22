package campaign

type Repository interface {
	Save(campaing *Campaign) error
}
