package vendor

import (
	"fmt"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/dbmodel"
	"github.com/jinzhu/gorm"
)

type dbRepository struct {
	db *gorm.DB
}

func FromMemorySQLite() (*dbRepository, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Not needed unless there are changes in model definition
	//dbmodel.MigrateModels(db)

	return &dbRepository{db: db}, nil
}

func FromPostgres(dbHost, dbPort, dbName, dbUser, dbPass string) (*dbRepository, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Not needed unless there are changes in model definition
	//dbmodel.MigrateModels(db)

	return &dbRepository{db: db}, nil
}

func (r *dbRepository) FindAll() ([]*Vendor, error) {
	vendors := make([]*dbmodel.Vendor, 0)
	r.db.Where("description != ''").Find(&vendors)
	return r.toVendors(vendors), nil
}

func (r *dbRepository) FindById(id string) (*Vendor, error) {
	dbVendor := &dbmodel.Vendor{}
	r.db.Where(&dbmodel.Vendor{Name: id}).First(&dbVendor)
	if dbVendor.ID == 0 {
		return nil, fmt.Errorf("not found")
	}
	return r.toVendor(dbVendor), nil
}

func (r *dbRepository) toVendors(dbVendors []*dbmodel.Vendor) []*Vendor {
	vendors := make([]*Vendor, 0)
	for _, dbVendor := range dbVendors {
		vendors = append(vendors, r.toVendor(dbVendor))
	}
	return vendors
}

func (r *dbRepository) toVendor(dbVendor *dbmodel.Vendor) *Vendor {
	vendor := &Vendor{}
	vendor.Kind = VENDOR
	vendor.Name = dbVendor.Name
	vendor.Description = dbVendor.Description
	vendor.Icon = dbVendor.Icon
	vendor.Website = dbVendor.Website
	return vendor
}
