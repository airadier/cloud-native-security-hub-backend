package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type fileRepository struct {
	path                     string
	resourcesCache           map[string][]*Resource
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
	resources = make([]*Resource, 0)

	for _, value := range f.resourcesCache {
		resources = append(resources, value...)
	}

	return resources, f.resourcesCacheError
}

func (f *fileRepository) FindAllLatestVersions() (resources []*Resource, err error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	resources = make([]*Resource, 0)

	for _, value := range f.resourcesCache {
		resources = append(resources, value[0])
	}

	return resources, f.resourcesCacheError
}

func (f *fileRepository) FindById(id string) ([]*Resource, error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	idToFind := strings.ToLower(id)

	if f.resourcesCacheError != nil {
		return nil, f.resourcesCacheError
	}

	if len(f.resourcesCache) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	if resources, ok := f.resourcesCache[idToFind]; ok {
		return resources, nil
	}

	return nil, fmt.Errorf("not found")
}

func (f *fileRepository) FindByIdLatestVersion(id string) (res *Resource, err error) {
	f.resourcesCacheFilledOnce.Do(f.fillResourcesCache)
	idToFind := strings.ToLower(id)

	if f.resourcesCacheError != nil {
		return nil, f.resourcesCacheError
	}

	if len(f.resourcesCache) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	if resources, ok := f.resourcesCache[idToFind]; ok {
		return resources[0], nil
	}

	return nil, fmt.Errorf("not found")
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

	if resources, ok := f.resourcesCache[idToFind]; ok {
		for _, resource := range resources {
			if resource.Version == version {
				return resource, nil
			}
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
	var resources = make(map[string][]*Resource)
	f.resourcesCacheError = filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			resource, err := resourceFromFile(path)
			if err != nil {
				return err
			}
			id := strings.ToLower(resource.ID)
			if val, ok := resources[id]; ok {
				resources[id] = append(val, &resource)
			} else {
				resources[id] = make([]*Resource, 1)
				resources[id][0] = &resource
			}
		}
		return nil
	})

	f.resourcesCache = resources

	//For each ID, sort by latest version first
	if f.resourcesCacheError == nil {
		for _, resources := range f.resourcesCache {
			sort.Sort(ByVersion(resources))
		}
	}
}
