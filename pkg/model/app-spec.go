package model

type AppSpec struct {
	Name            string               `json:"name" yaml:"name"`
	Image           string               `json:"image" yaml:"image"`
	Env             []Env                `json:"env" yaml:"env"`
	Secrets         *SecretSpec          `json:"secrets" yaml:"secrets"`
	Sidecar         []string             `json:"sidecar" yaml:"sidecar"`
	Service         *ServiceSpec         `json:"service" yaml:"service"`
	Workload        *WorkloadSpec        `json:"workload" yaml:"workload"`
	ServiceAccount  *string              `json:"serviceAccount" yaml:"serviceAccount"`
	Resources       []string             `json:"resources" yaml:"resources"`
	SecurityContext *SecurityContextSpec `json:"securityContext" yaml:"securityContext"`
	Mixins          []string             `json:"mixins" yaml:"mixins"`
	Ingress         *IngressSpec         `json:"ingress" yaml:"ingress"`
	Mounts          []string             `json:"mounts" yaml:"mounts"`
	Args            *ArgsSpec             `json:"args" yaml:"args"`
}

type Env struct {
	EnvSet *string `json:"env-set" yaml:"env-set"`
	Name   *string `json:"name" yaml:"name"`
	Value  *string `json:"value" yaml:"value"`
}

type SidecarSpec struct {
	Name     string `json:"name" yaml:"name"`
	Image    string `json:"image" yaml:"image"`
	Template string `json:"template" yaml:"template"`
}

type SecretSpec struct {
	Enabled  bool    `json:"enabled" yaml:"enabled"`
	Strategy *string `json:"strategy" yaml:"strategy"`
}

type ServiceSpec struct {
	Port           map[string]int `json:"port" yaml:"port"`
	HealthCheckUrl *string        `json:"healthCheck" yaml:"healthCheck"`
	Headless       *bool          `json:"headless" yaml:"headless"`
}

type WorkloadSpec struct {
	Target  string  `json:"target" yaml:"target"`
	Scaling *string `json:"scaling" yaml:"scaling"`
}

type SecurityContextSpec struct {
	RunAsUser                string `json:"runAsUser" yaml:"runAsUser"`
	AllowPrivilegeEscalation *bool  `json:"allowPrivilegeEscalation" yaml:"allowPrivilegeEscalation"`
	ReadOnlyRootFilesystem   *bool  `json:"readOnlyRootFilesystem" yaml:"readOnlyRootFilesystem"`
	RunAsNonRoot             *bool  `json:"runAsNonRoot" yaml:"runAsNonRoot"`
}

type IngressSpec struct {
	Name  string `json:"name" yaml:"name"`
	Group string `json:"group" yaml:"group"`
}
