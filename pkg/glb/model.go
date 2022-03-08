package glb

type Environment struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Spec EnvironmentSpec `json:"spec" yaml:"spec"`
}
type Metadata struct {
	Name string `json:"name" yaml:"name"`
}
type EnvironmentSpec struct {
	Selector map[string]string `json:"selector" yaml:"selector"`
	Data []KeyValue `json:"data" yaml:"data"`
}

type KeyValue struct {
	Name string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

