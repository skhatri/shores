package mixin

import (
	"github.com/skhatri/shores/pkg/functions"
	"github.com/skhatri/shores/pkg/model"
)

func LoadMixins(files []string) map[string]model.MixinTemplate {
	errors := make([]string, 0)
	mixins := make(map[string]model.MixinTemplate, 0)
	for _, file := range files {
		mixinKind := Mixin{}
		err := functions.UnmarshalFile(file, &mixinKind)
		if err != nil {
			errors = append(errors, err.Error())
		}
		if mixinKind.Kind != "Mixin" {
			continue
		}
		mixins[mixinKind.Metadata.Name] = mixinKind.Spec.Template
	}
	return mixins
}
