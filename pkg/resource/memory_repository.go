package resource

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

type MemoryRepository struct {
	resources []*Resource
}

func NewMemoryRepository(resources []*Resource) Repository {
	return &MemoryRepository{
		resources: resources,
	}
}

func (r *MemoryRepository) FindAll() ([]*Resource, error) {
	return r.resources, nil
}

func (r *MemoryRepository) FindById(id string) (*Resource, error) {
	version, err := r.getLatestVersionForId(id)
	if err != nil {
		return nil, err
	}
	return r.FindByIdAndVersion(id, version)
}

func (r *MemoryRepository) getLatestVersionForId(id string) (version string, err error) {
	idToFind := strings.ToLower(id)

	var latestVersion *semver.Version
	for _, resource := range r.resources {
		if resource.ID == idToFind {
			resourceVersion, err := semver.NewVersion(resource.Version)

			if err == nil && (latestVersion == nil || resourceVersion.GreaterThan(latestVersion)) {
				latestVersion = resourceVersion
			}
		}
	}

	if latestVersion == nil {
		return "", fmt.Errorf("not found")
	}
	return latestVersion.String(), nil
}

func (r *MemoryRepository) FindByIdAndVersion(id, version string) (*Resource, error) {
	idToFind := strings.ToLower(id)
	for _, res := range r.resources {
		if res.ID == idToFind && res.Version == version {
			return res, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *MemoryRepository) Add(resource Resource) {
	r.resources = append(r.resources, &resource)
}
