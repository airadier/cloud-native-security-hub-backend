package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/dbmodel"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/vendor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	checkArgs()

	dbHost, dbPort, dbName, dbUser, dbPass := checkDbParameters()
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)
	db, err := gorm.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbmodel.MigrateModels(db)

	if db.Error != nil {
		log.Fatal(err)
	}

	importVendors(db, os.Args[1])

	importResources(db, os.Args[2])
}

func checkArgs() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <path_to_vendor_yamls> <path_to_resource_yamls>\n", os.Args[0])
		os.Exit(-1)
	}
}

func checkDbParameters() (dbHost, dbPort, dbName, dbUser, dbPass string) {
	dbHost, hostOk := os.LookupEnv("DB_HOST")
	if !hostOk {
		log.Fatal("The DB_HOST env var is not set")
	}

	dbPort, portOk := os.LookupEnv("DB_PORT")
	if !portOk {
		log.Fatal("The DB_PORT env var is not set")
	}

	dbName, nameOk := os.LookupEnv("DB_NAME")
	if !nameOk {
		log.Fatal("The DB_NAME env var is not set")
	}

	dbUser, userOk := os.LookupEnv("DB_USER")
	if !userOk {
		log.Fatal("The DB_USER env var is not set")
	}

	dbPass, passOk := os.LookupEnv("DB_PASS")
	if !passOk {
		log.Fatal("The DB_PASS env var is not set")
	}

	return dbHost, dbPort, dbName, dbUser, dbPass
}

func importVendors(db *gorm.DB, vendorsPath string) {
	vendorsRepo, err := vendor.FromPath(vendorsPath)
	if err != nil {
		log.Fatal(err)
	}

	vendors, err := vendorsRepo.FindAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, vendor := range vendors {
		dbVendor := &dbmodel.Vendor{}

		db.Where(&dbmodel.Vendor{Name: vendor.Name}).First(dbVendor)

		dbVendor.Name = vendor.Name
		dbVendor.Description = vendor.Description
		dbVendor.Icon = vendor.Icon
		dbVendor.Website = vendor.Website
		dbVendor.Resources = make([]dbmodel.Resource, 0)

		if dbVendor.ID == 0 {
			fmt.Printf("Importing vendor %s\n", vendor.Name)

			if err := db.Create(dbVendor).Error; err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Updating vendor %s\n", vendor.Name)
		}
	}
}

func importResources(db *gorm.DB, resourcesPath string) {
	resourcesRepo, err := resource.FromPath(resourcesPath)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := resourcesRepo.FindAllLatestVersions()
	if err != nil {
		log.Fatal(err)
	}

	for _, resource := range resources {
		fmt.Printf("Importing resource %s\n", resource.Name)

		existingRes := &dbmodel.Resource{}
		db.Where(&dbmodel.Resource{ResourceID: resource.ID}).First(existingRes)

		if existingRes.ID != 0 {
			fmt.Printf("Already exists, skipping\n")
			continue
		}

		res, err := dbResourceFromID(resourcesRepo, resource.ID)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Create(res).Error; err != nil {
			log.Fatal(err)
		}

		if resource.Vendor != "" {

			vendor := &dbmodel.Vendor{}
			db.Where(&dbmodel.Vendor{Name: resource.Vendor}).First(vendor)

			if vendor.ID == 0 {
				vendor.Name = resource.Vendor
				vendor.Resources = []dbmodel.Resource{*res}
				fmt.Printf("  New vendor %s\n", resource.Vendor)
				db.Create(vendor)
			} else {
				fmt.Printf("  Adding to vendor %s\n", resource.Vendor)
				db.Model(vendor).Association("Resources").Append(res)
			}
		}
	}
}

func dbResourceFromID(r resource.Repository, resourceId string) (*dbmodel.Resource, error) {
	dbr := &dbmodel.Resource{}
	dbr.ResourceID = resourceId
	dbr.Versions = make([]dbmodel.ResourceVersion, 0)
	resources, err := r.FindById(resourceId)
	if err != nil {
		return dbr, err
	}
	for _, res := range resources {
		dbr.Versions = append(dbr.Versions, fromResource(res))
	}

	return dbr, nil
}

func fromResource(r *resource.Resource) dbmodel.ResourceVersion {
	version := dbmodel.ResourceVersion{}
	version.Name = r.Name
	version.ShortDescription = r.ShortDescription
	version.Description = r.Description
	version.Keywords = strings.Join(r.Keywords, ",")
	version.Icon = r.Icon
	version.Website = r.Website
	version.Version = r.Version

	for _, m := range r.Maintainers {

		version.Maintainers = append(version.Maintainers, dbmodel.Maintainer{
			Name:  m.Name,
			Email: m.Email,
		})
	}

	for _, v := range r.Rules {
		version.Rules = append(version.Rules, dbmodel.FalcoRuleData{
			Raw: v.Raw,
		})
	}

	return version
}
