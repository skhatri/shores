package model

import (
	"fmt"
	"github.com/skhatri/shores/pkg/functions"
	"strings"
)

type ProductSet struct {
	Namespace          *string        `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	ContainerNamespace *string        `json:"containerNamespace,omitempty" yaml:"container_namespace,omitempty"`
	Apps               []*ReleaseSpec `json:"apps,omitempty" yaml:"apps,omitempty"`
}

type ReleaseSpec struct {
	Name      string  `json:"name" yaml:"name"`
	Image     *string `json:"image,omitempty" yaml:"image,omitempty"`
	Version   *string `json:"version,omitempty" yaml:"version,omitempty"`
	Namespace string  `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

func NewProductSetFromFile(file string, namespaceOverride string) (*ProductSet, error) {
	var productSet ProductSet
	err := functions.UnmarshalFile(file, &productSet)
	if err != nil {
		return nil, err
	}
	namespace := "default"
	if namespaceOverride != "" {
		productSet.Namespace = &namespaceOverride
	}
	if productSet.Namespace == nil {
		productSet.Namespace = &namespace
	}
	prefix := ""
	if productSet.ContainerNamespace != nil {
		prefix = fmt.Sprintf("%s/", *productSet.ContainerNamespace)
	}
	for _, appRef := range productSet.Apps {
		appRef.Namespace = *productSet.Namespace
		if appRef.Image == nil {
			version := "latest"
			if appRef.Version != nil {
				version = *appRef.Version
			}
			imageNameFromVersion := fmt.Sprintf("%s:%s", appRef.Name, version)
			appRef.Image = &imageNameFromVersion
		}
		if prefix != "" {
			imageWithPrefix := fmt.Sprintf("%s%s", prefix, *appRef.Image)
			appRef.Image = &imageWithPrefix
		}
		if appRef.Version == nil {
			parts := strings.Split(*appRef.Image, ":")
			if len(parts) > 1 {
				appRef.Version = &parts[1]
			}
		}
	}
	return &productSet, nil
}

type Task struct {
	User      string `json:"user" yaml:"user"`
	ReleaseId string `json:"releaseId" yaml:"releaseId"`
	Action    string `json:"action" yaml:"action"`
	Command   string `json:"command" yaml:"command"`
	Created   string `json:"created" yaml:"created"`
	ChangeRef string `json:"changeRef" yaml:"changeRef"`
}
