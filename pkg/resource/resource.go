package resource

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type Kind string

const (
	FALCO_RULE Kind = "FalcoRule"
)

type Resource struct {
	ID               string           `json:"id,omitempty" yaml:"id,omitempty"`
	Kind             Kind             `json:"kind" yaml:"kind"`
	Vendor           string           `json:"vendor" yaml:"vendor"`
	Name             string           `json:"name" yaml:"name"`
	ShortDescription string           `json:"shortDescription" yaml:"shortDescription"`
	Description      string           `json:"description" yaml:"description"`
	Keywords         []string         `json:"keywords" yaml:"keywords"`
	Icon             string           `json:"icon" yaml:"icon"`
	Website          string           `json:"website" yaml:"website"`
	Maintainers      []*Maintainer    `json:"maintainers" yaml:"maintainers"`
	LatestVersion    string           `json:"latestVersion" yaml:"latestVersion"`
	Rules            []*FalcoRuleData `json:"rules" yaml:"rules"`
}

func (r *Resource) GenerateRulesForHelmChart(version string) ([]byte, error) {
	raw := make(map[string]map[string]string)
	raw["customRules"] = map[string]string{}

	version = r.getLatestVersion(version)

	found := false
	for _, rule := range r.Rules {
		if rule.Version == version {
			raw["customRules"]["rules-"+r.ID+".yaml"] += rule.Raw
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("not found")
	}

	result, _ := yaml.Marshal(raw)
	return result, nil
}

func (r *Resource) getLatestVersion(requestedVersion string) string {
	if requestedVersion != "" {
		return requestedVersion
	}

	if r.LatestVersion != "" {
		// If version not specified, take from resource.latestVersion
		return r.LatestVersion
	}

	// As fallback, take version from last defined rule
	return r.Rules[len(r.Rules)-1].Version
}

type resourceAlias Resource // Avoid stack overflow while marshalling / unmarshalling

func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	res := resourceAlias{}
	err = unmarshal(&res)
	if err != nil {
		return
	}
	*r = Resource(res)
	r.ID = r.generateID()
	return
}

func (r *Resource) MarshalYAML() (interface{}, error) {
	x := resourceAlias(*r)
	x.ID = r.generateID()
	return yaml.Marshal(x)
}

func (r *Resource) UnmarshalJSON(data []byte) (err error) {
	res := resourceAlias{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	*r = Resource(res)
	r.ID = r.generateID()
	return
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	x := resourceAlias(*r)
	x.ID = r.generateID()
	return json.Marshal(x)
}

type Maintainer struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
}

type FalcoRuleData struct {
	Version            string `json:"version" yaml:"version"`
	TargetFalcoVersion string `json:"targetFalcoVersion" yaml:"targetFalcoVersion"`
	Description        string `json:"description" yaml:"description"`
	Raw                string `json:"raw" yaml:"raw"`
}

func (r *Resource) Validate() error {
	var errors []string

	if r.Kind == "" {
		errors = append(errors, "the resource must have a defined Kind")
	}
	if r.Vendor == "" {
		errors = append(errors, "the resource must be assigned to a vendor")
	}
	if len(r.Maintainers) == 0 {
		errors = append(errors, "the resource must have at least one maintainer")
	}
	if r.Icon == "" {
		errors = append(errors, "the resource must have a valid icon")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ","))
	}

	return nil
}

func (r *Resource) generateID() string {
	return strings.ToLower(r.Name)
}
