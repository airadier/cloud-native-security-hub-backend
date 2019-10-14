package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/dbmodel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type dbRepository struct {
	db *gorm.DB
}

func FromMemorySQLite() (*dbRepository, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	dbmodel.MigrateModels(db)

	return &dbRepository{db: db}, nil
}

func FromPostgres(dbHost, dbPort, dbName, dbUser, dbPass string) (*dbRepository, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	dbmodel.MigrateModels(db)

	return &dbRepository{db: db}, nil
}

func (r *dbRepository) FindAll() ([]*Resource, error) {
	resources := make([]*dbmodel.Resource, 0)
	r.db.Find(&resources)
	return r.toResources(resources), nil
}

func (r *dbRepository) FindAllLatestVersions() ([]*Resource, error) {
	resources := make([]*dbmodel.Resource, 0)
	r.db.Find(&resources)
	return r.toResourcesLatestVersion(resources), nil
}

func (r *dbRepository) FindById(id string) ([]*Resource, error) {
	resources := make([]*dbmodel.Resource, 0)
	r.db.Where(&dbmodel.Resource{ResourceID: id}).Find(&resources)
	if len(resources) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return r.toResources(resources), nil
}

func (r *dbRepository) FindByIdLatestVersion(id string) (res *Resource, err error) {
	resource := &dbmodel.Resource{}
	r.db.Where(&dbmodel.Resource{ResourceID: id}).First(resource)
	if resource.ID == 0 {
		return nil, fmt.Errorf("not found")
	}
	return r.toResourcesLatestVersion([]*dbmodel.Resource{resource})[0], nil
}

func (r *dbRepository) FindByIdAndVersion(id, version string) (res *Resource, err error) {
	dbResource := &dbmodel.Resource{}
	dbResourceVersion := &dbmodel.ResourceVersion{}
	dbVendor := &dbmodel.Vendor{}

	r.db.Where(&dbmodel.Resource{ResourceID: id}).First(dbResource)
	if dbResource.ID == 0 {
		return nil, fmt.Errorf("not found")
	}

	r.db.Where(&dbmodel.ResourceVersion{ResourceID: dbResource.ID, Version: version}).First(dbResourceVersion)
	if dbResourceVersion.ID == 0 {
		return nil, fmt.Errorf("not found")
	}

	r.db.Where(&dbmodel.Vendor{ID: *dbResource.VendorID}).First(dbVendor)
	res = r.toResourceVersion(dbResourceVersion)
	res.ID = dbResource.ResourceID
	res.Vendor = dbVendor.Name
	return res, nil
}

func (r *dbRepository) toResources(dbResources []*dbmodel.Resource) []*Resource {
	resources := make([]*Resource, 0)
	for _, dbRes := range dbResources {
		r.db.Model(dbRes).Association("Versions").Find(&dbRes.Versions)
		resources = append(resources, r.toResourceVersions(dbRes)...)
	}
	return resources
}

func (r *dbRepository) toResourcesLatestVersion(dbResources []*dbmodel.Resource) []*Resource {
	resources := make([]*Resource, 0)
	for _, dbRes := range dbResources {
		r.db.Model(dbRes).Association("Versions").Find(&dbRes.Versions)

		versions := r.toResourceVersions(dbRes)
		sort.Sort(ByVersion(versions))

		resources = append(resources, versions[0])
	}
	return resources
}

func (r *dbRepository) toResourceVersions(dbRes *dbmodel.Resource) []*Resource {

	dbVendor := &dbmodel.Vendor{}
	if dbRes.VendorID != nil {
		r.db.Where(&dbmodel.Vendor{ID: *dbRes.VendorID}).First(dbVendor)
	}

	resources := make([]*Resource, 0)
	for _, dbVersion := range dbRes.Versions {
		resourceVersion := r.toResourceVersion(&dbVersion)
		resourceVersion.Vendor = dbVendor.Name
		resourceVersion.ID = dbRes.ResourceID
		resources = append(resources, resourceVersion)
	}
	return resources
}

func (r *dbRepository) toResourceVersion(dbRes *dbmodel.ResourceVersion) *Resource {
	resource := &Resource{}
	resource.Kind = FALCO_RULE
	resource.Name = dbRes.Name
	resource.ShortDescription = dbRes.ShortDescription
	resource.Description = dbRes.Description
	resource.Keywords = strings.Split(dbRes.Keywords, ",")
	resource.Icon = dbRes.Icon
	resource.Website = dbRes.Website
	resource.Maintainers = r.toMaintainers(dbRes)
	resource.Version = dbRes.Version
	resource.Rules = r.ToFalcoRules(dbRes)
	return resource
}

func (r *dbRepository) toMaintainers(dbRes *dbmodel.ResourceVersion) []*Maintainer {
	maintainers := make([]*Maintainer, 0)
	r.db.Model(dbRes).Association("Maintainers").Find(&dbRes.Maintainers)
	for _, maintainer := range dbRes.Maintainers {
		maintainers = append(maintainers, &Maintainer{Name: maintainer.Name, Email: maintainer.Email})
	}
	return maintainers
}

func (r *dbRepository) ToFalcoRules(dbRes *dbmodel.ResourceVersion) []*FalcoRuleData {
	rules := make([]*FalcoRuleData, 0)
	r.db.Model(dbRes).Association("Rules").Find(&dbRes.Rules)
	for _, rule := range dbRes.Rules {
		rules = append(rules, &FalcoRuleData{Raw: rule.Raw})
	}
	return rules
}
