package resource

import (
	"github.com/skhatri/shores/pkg/functions"
	"github.com/skhatri/shores/pkg/model"
)

type ResourceKind struct {
	Kind     string      `json:"kind" yaml:"kind"`
	Metadata Metadata    `json:"metadata" yaml:"metadata"`
	Spec     ResourceDef `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type ResourceDef struct {
	Selector map[string]string `json:"selector" yaml:"selector"`
	Data     model.Resources   `json:"data" yaml:"data"`
}

func LoadResources(files []string) map[string]model.Resources {
	errors := make([]error, 0)
	resources := make(map[string]model.Resources, 0)
	for _, file := range files {
		resourceKind := ResourceKind{}
		err := functions.UnmarshalFile(file, &resourceKind)
		if err != nil {
			errors = append(errors, err)
		}
		if resourceKind.Kind != "Resource" {
			continue
		}
		resources[resourceKind.Metadata.Name] = resourceKind.Spec.Data
	}
	return resources
}
