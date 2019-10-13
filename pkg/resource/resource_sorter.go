package resource

import "github.com/Masterminds/semver"

// ByVersion implements sort.Interface based on the Version field.
type ByVersion []*Resource

func (r ByVersion) Len() int {
	return len(r)
}

func (r ByVersion) Less(i, j int) bool {

	v1, err1 := semver.NewVersion(r[i].Version)
	v2, err2 := semver.NewVersion(r[j].Version)

	if err1 != nil || err2 != nil {
		return false
	}

	// Higher versions get first in the list
	return v1.GreaterThan(v2)
}

func (r ByVersion) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
