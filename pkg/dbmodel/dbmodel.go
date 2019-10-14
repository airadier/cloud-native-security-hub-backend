package dbmodel

import (
	"log"

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
