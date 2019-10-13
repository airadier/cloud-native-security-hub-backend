package resource

import (
	"fmt"
	"sort"
	"strings"
)

type MemoryRepository struct {
	resources map[string][]*Resource
}

// NewMemoryRepository created a new memory repository.
// Please provide resources grouped by ID and sorted newer versions forst
func NewMemoryRepository(resources []*Resource) Repository {
	r := &MemoryRepository{resources: make(map[string][]*Resource)}
	for _, res := range resources {
		r.Add(*res)
	}

	return r
}

func (r *MemoryRepository) FindAll() ([]*Resource, error) {
	resources := make([]*Resource, 0)
	for _, res := range r.resources {
		resources = append(resources, res...)
	}

	return resources, nil
}

func (r *MemoryRepository) FindAllLatestVersions() ([]*Resource, error) {
	resources := make([]*Resource, 0)
	for _, res := range r.resources {
		resources = append(resources, res[0])
	}

	return resources, nil
}

func (r *MemoryRepository) FindById(id string) ([]*Resource, error) {
	idToFind := strings.ToLower(id)

	if resources, ok := r.resources[idToFind]; ok {
		return resources, nil
	}

	return nil, fmt.Errorf("not found")
}

func (r *MemoryRepository) FindByIdLatestVersion(id string) (*Resource, error) {
	idToFind := strings.ToLower(id)

	if resources, ok := r.resources[idToFind]; ok {
		return resources[0], nil
	}

	return nil, fmt.Errorf("not found")
}

func (r *MemoryRepository) FindByIdAndVersion(id, version string) (*Resource, error) {
	idToFind := strings.ToLower(id)

	if resources, ok := r.resources[idToFind]; ok {
		for _, res := range resources {
			if res.Version == version {
				return res, nil
			}
		}
	}

	return nil, fmt.Errorf("not found")
}

func (r *MemoryRepository) Add(resource Resource) {
	idToFind := strings.ToLower(resource.ID)

	if resourceList, ok := r.resources[idToFind]; ok {
		r.resources[idToFind] = append(resourceList, &resource)
	} else {
		r.resources[idToFind] = make([]*Resource, 1)
		r.resources[idToFind][0] = &resource
	}

	sort.Sort(ByVersion(r.resources[idToFind]))
}
