package resource

type Repository interface {
	FindAll() ([]*Resource, error)
	FindById(id string) (*Resource, error)
	FindByIdAndVersion(id, version string) (*Resource, error)
}
