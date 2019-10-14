package dbmodel

import (
	"log"
	"strings"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/jinzhu/gorm"
)

type Vendor struct {
	ID          uint
	Name        string `gorm:"unique_index"`
	Description string
	Icon        string
	Website     string
	Resources   []Resource `gorm:"foreignkey:VendorID"`
}

type Resource struct {
	ID         uint
	ResourceID string `gorm:"unique_index"`
	Versions   []ResourceVersion
	VendorID   *uint
}

type ResourceVersion struct {
	ID               uint
	Name             string
	ShortDescription string
	Description      string
	Keywords         string
	Icon             string
	Website          string
	Maintainers      []Maintainer
	Version          string
	Rules            []FalcoRuleData
	ResourceID       uint
}

type Maintainer struct {
	ID                uint
	Name              string
	Email             string
	ResourceVersionID uint
}

type FalcoRuleData struct {
	ID                uint
	Raw               string
	ResourceVersionID uint
}

func MigrateModels(db *gorm.DB) {
	if err := db.AutoMigrate(
		&FalcoRuleData{},
		&Maintainer{},
		&ResourceVersion{},
		&Resource{},
		&Vendor{}).Error; err != nil {
		log.Fatal(err)
	}

	db.Model(&Resource{}).AddForeignKey("vendor_id", "vendors(id)", "RESTRICT", "RESTRICT")
	db.Model(&ResourceVersion{}).AddForeignKey("resource_id", "resources(id)", "RESTRICT", "RESTRICT")
	db.Model(&Maintainer{}).AddForeignKey("resource_version_id", "resource_versions(id)", "RESTRICT", "RESTRICT")
	db.Model(&FalcoRuleData{}).AddForeignKey("resource_version_id", "resource_versions(id)", "RESTRICT", "RESTRICT")
}

func FromResourceID(r resource.Repository, resourceId string) (*Resource, error) {
	dbr := &Resource{}
	dbr.ResourceID = resourceId
	dbr.Versions = make([]ResourceVersion, 0)
	resources, err := r.FindById(resourceId)
	if err != nil {
		return dbr, err
	}
	for _, res := range resources {
		dbr.Versions = append(dbr.Versions, FromResource(res))
	}

	return dbr, nil
}

func FromResource(r *resource.Resource) ResourceVersion {
	version := ResourceVersion{}
	version.Name = r.Name
	version.ShortDescription = r.ShortDescription
	version.Description = r.Description
	version.Keywords = strings.Join(r.Keywords, ",")
	version.Icon = r.Icon
	version.Website = r.Website
	version.Version = r.Version

	for _, m := range r.Maintainers {

		version.Maintainers = append(version.Maintainers, Maintainer{
			Name:  m.Name,
			Email: m.Email,
		})
	}

	for _, v := range r.Rules {
		version.Rules = append(version.Rules, FalcoRuleData{
			Raw: v.Raw,
		})
	}

	return version
}
