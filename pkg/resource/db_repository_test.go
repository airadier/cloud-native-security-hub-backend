package resource

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/dbmodel"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	repo, err := FromMemorySQLite()

	assert.NoError(t, err)
	defer repo.db.Close()

	err = initSampleData(repo.db)
	assert.NoError(t, err)

	resources, err := repo.FindAll()
	assert.NoError(t, err)

	assert.Equal(t, 4, len(resources))
}

func TestFindLatestVersions(t *testing.T) {
	repo, err := FromMemorySQLite()
	assert.NoError(t, err)
	defer repo.db.Close()

	err = initSampleData(repo.db)
	assert.NoError(t, err)

	resources, err := repo.FindAllLatestVersions()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resources))
}

func initSampleData(db *gorm.DB) error {
	if err := db.Create(&dbmodel.Resource{
		ResourceID: "apache",
		Versions: []dbmodel.ResourceVersion{
			{
				Name:        "apache v1",
				Version:     "1.0.0",
				Description: "blahblahv1",
				Maintainers: []dbmodel.Maintainer{
					{Name: "Maintainer1", Email: "email1"},
					{Name: "Maintainer2", Email: "email2"},
				},
				Rules: []dbmodel.FalcoRuleData{
					{Raw: "rawData1"},
				},
			},
			{
				Name:        "apache v2",
				Version:     "2.0.0",
				Description: "blahblahv2",
				Maintainers: []dbmodel.Maintainer{
					{Name: "Maintainer1", Email: "email1"},
					{Name: "Maintainer2", Email: "email2"},
				},
				Rules: []dbmodel.FalcoRuleData{
					{Raw: "rawData2"},
				},
			},
			{
				Name:        "apache v1.1",
				Version:     "1.1.0",
				Description: "blahblahv1.1",
				Maintainers: []dbmodel.Maintainer{
					{Name: "Maintainer1", Email: "email1"},
					{Name: "Maintainer2", Email: "email2"},
				},
				Rules: []dbmodel.FalcoRuleData{
					{Raw: "rawData1.1"},
				},
			},
		},
	}).Error; err != nil {
		return err
	}

	if err := db.Create(&dbmodel.Resource{
		ResourceID: "mongo",
		Versions: []dbmodel.ResourceVersion{
			{
				Name:        "mongo v1",
				Version:     "1.0.0",
				Description: "blahblahv1",
				Maintainers: []dbmodel.Maintainer{
					{Name: "Maintainer1", Email: "email1"},
					{Name: "Maintainer2", Email: "email2"},
				},
				Rules: []dbmodel.FalcoRuleData{
					{Raw: "rawData1"},
				},
			},
		},
	}).Error; err != nil {
		return err
	}

	return nil
}
