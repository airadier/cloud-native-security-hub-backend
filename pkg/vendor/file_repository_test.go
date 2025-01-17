package vendor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileRepositoryWalksADirectoryAndExtractVendors(t *testing.T) {
	path := "../../test/fixtures/vendors"
	vendorRepository, _ := FromPath(path)

	vendors, _ := vendorRepository.FindAll()

	assert.Equal(t, buildVendorsFromFixtures(), vendors)
}

func buildVendorsFromFixtures() []*Vendor {
	return []*Vendor{
		{
			ID:          "apache",
			Kind:        "Vendor",
			Name:        "Apache",
			Description: "# Apache Software Foundation\n",
			Icon:        "https://upload.wikimedia.org/wikipedia/commons/thumb/d/db/Apache_Software_Foundation_Logo_%282016%29.svg/2560px-Apache_Software_Foundation_Logo_%282016%29.svg.png",
			Website:     "https://apache.org/",
		},

		{
			ID:          "mongo",
			Kind:        "Vendor",
			Name:        "Mongo",
			Description: "# MongoDB Inc.\n",
			Icon:        "https://upload.wikimedia.org/wikipedia/en/thumb/4/45/MongoDB-Logo.svg/640px-MongoDB-Logo.svg.png",
			Website:     "https://mongodb.com/",
		},
	}
}

func TestFileRepositoryReturnsAnErrorIfPathDoesNotExist(t *testing.T) {
	nonExistentPath := "../foo"

	_, err := FromPath(nonExistentPath)

	assert.Error(t, err)
}
