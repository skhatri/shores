package mixin

import "github.com/skhatri/shores/pkg/model"

type Mixin struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Spec MixinSpec `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type MixinSpec struct {
	Selector map[string]string `json:"selector" yaml:"selector"`
	Template model.MixinTemplate `json:"template" yaml:"template"`
}
