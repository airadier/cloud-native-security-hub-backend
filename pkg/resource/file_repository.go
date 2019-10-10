package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

type fileRepository struct {
	path                     string
	resourcesCache           []*Resource
	resourcesCacheFilledOnce sync.Once
	resourcesCacheError      error
}

func FromPath(path string) (*fileRepository, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return &fileRepository{path: path}, nil
}

func (f *fileRepository) FindAll() (resources []*Resource, err error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	return f.resourcesCache, f.resourcesCacheError
}

func (f *fileRepository) FindById(id string) (res *Resource, err error) {
	latestVersion, err := f.getLatestVersionForId(id)
	if err != nil {
		return nil, err
	}

	return f.FindByIdAndVersion(id, latestVersion)
}

func (f *fileRepository) getLatestVersionForId(id string) (version string, err error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	idToFind := strings.ToLower(id)

	if f.resourcesCacheError != nil {
		return "", f.resourcesCacheError
	}

	if len(f.resourcesCache) == 0 {
		return "", fmt.Errorf("no resources")
	}

	var latestVersion *semver.Version
	for _, resource := range f.resourcesCache {
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

func (f *fileRepository) FindByIdAndVersion(id, version string) (res *Resource, err error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	idToFind := strings.ToLower(id)

	if f.resourcesCacheError != nil {
		return nil, f.resourcesCacheError
	}

	if len(f.resourcesCache) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	for _, resource := range f.resourcesCache {
		if resource.ID == idToFind && resource.Version == version {
			res = resource
			return
		}
	}

	err = fmt.Errorf("not found")
	return
}

func resourceFromFile(path string) (resource Resource, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	err = yaml.NewDecoder(file).Decode(&resource)
	if err != nil {
		return
	}

	return
}

func (f *fileRepository) fillResourcesCache() {
	var resources []*Resource
	f.resourcesCacheError = filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			resource, err := resourceFromFile(path)
			if err != nil {
				return err
			}
			resources = append(resources, &resource)
		}
		return nil
	})
	f.resourcesCache = resources
}
