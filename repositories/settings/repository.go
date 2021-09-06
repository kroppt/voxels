package settings

// Repository stores user settings.
type Repository struct {
	c core
}

// New creates and returns a new settings repository.
func New(fileMod fileMod) *Repository {
	return &Repository{
		core{
			fileMod: fileMod,
		},
	}
}
