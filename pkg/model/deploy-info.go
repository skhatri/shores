package model

type DeploymentSummary struct {
	Namespace string
	Items     []DeploymentItem
}

type DeploymentItem struct {
	Name string
	Kind string
	Path string
}
