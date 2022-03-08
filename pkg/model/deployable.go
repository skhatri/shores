package model

type Deployable struct {
	Kind               string               `json:"kind"`
	Namespace          string               `json:"namespace"`
	Artifact           ArtifactInfo         `json:"artifact"`
	Checks             *Healthcheck         `json:"checks"`
	Target             TargetInfo           `json:"target"`
	Env                map[string]string    `json:"env"`
	InitContainer      []InitContainerInfo  `json:"initContainer"`
	Sidecar            []SidecarInfo        `json:"sidecar"`
	Service            []ServiceInfo        `json:"service"`
	Metadata           Metadata             `json:"metadata"`
	ServiceAccountName *string              `json:"serviceAccountName"`
	Runtime            RuntimeArgs          `json:"runtime"`
	ServiceEnabled     bool                 `json:"serviceEnabled"`
	Resources          *Resources           `json:"resources"`
	SecurityContext    *SecurityContextSpec `json:"securityContext"`
	Ingress            *IngressSpec         `json:"ingress"`
}

func (d *Deployable) Ports() []PortType {
	ports := make([]PortType, 0)
	svc := d.Service
	containerPorts := make(map[int]struct{}, 0)
	for _, s := range svc {
		for _, p := range s.Port {
			_, ok := containerPorts[p.Port]
			if !ok {
				ports = append(ports, p)
				containerPorts[p.Port] = struct{}{}
			}
		}
	}
	return ports
}

type Metadata struct {
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
}

type ArtifactInfo struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}

type ReleaseInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Service string `json:"service,omitempty"`
}

type InitContainerInfo struct{}

type SidecarInfo struct{}

type ServiceInfo struct {
	Headless bool       `json:"headless,omitempty"`
	Port     []PortType `json:"port"`
}

type PortType struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port,omitempty"`
	TargetPort string `json:"targetPort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}

type Healthcheck struct {
	Port string `json:"port"`
	Path string `json:"path"`
}

type TargetInfo struct {
	Replica      int               `json:"replica"`
	NodeSelector map[string]string `json:"nodeSelector"`
}

type RuntimeArgs struct {
	Entrypoint string   `json:"entrypoint"`
	Command    []string `json:"command"`
}

type Resources struct {
	Requests *ResourceValue `json:"requests"`
	Limits   *ResourceValue `json:"limits"`
}

type ResourceValue struct {
	Cpu    *string `json:"cpu"`
	Memory *string `json:"memory"`
}
