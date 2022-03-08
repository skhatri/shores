package model

type MixinTemplate struct {
	Secrets         *SecretSpec          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Sidecar         []*SidecarSpec       `json:"sidecar,omitempty" yaml:"sidecar,omitempty"`
	Service         *ServiceSpec         `json:"service,omitempty" yaml:"service,omitempty"`
	Workload        *WorkloadSpec        `json:"workload,omitempty" yaml:"workload,omitempty"`
	Resources       []*string            `json:"resources,omitempty" yaml:"resources,omitempty"`
	SecurityContext *SecurityContextSpec `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`
	Args            *ArgsSpec             `json:"args,omitempty" yaml:"args,omitempty"`
}

func (mx *MixinTemplate) Merge(other *MixinTemplate) *MixinTemplate {
	newTemplate := MixinTemplate{}

	myWorkload := mx.Workload
	theirWorkload := other.Workload
	newTemplate.Workload = myWorkload
	if theirWorkload != nil {
		newTemplate.Workload = theirWorkload
	}

	mergedResources := make([]*string, 0)
	myResources := mx.Resources
	theirResources := other.Resources
	if len(myResources) > 0 {
		for _, resourceName := range myResources {
			mergedResources = append(mergedResources, resourceName)
		}
	}
	if len(theirResources) > 0 {
		for _, resourceName := range theirResources {
			mergedResources = append(mergedResources, resourceName)
		}
	}
	newTemplate.Resources = mergedResources

	mySecrets := mx.Secrets
	theirSecrets := other.Secrets
	newTemplate.Secrets = mySecrets
	if theirSecrets != nil {
		newTemplate.Secrets = theirSecrets
	}

	myService := mx.Service
	theirService := other.Service
	newTemplate.Service = myService
	if theirService != nil {
		newTemplate.Service = theirService
	}

	mySecurityContext := mx.SecurityContext
	theirSecurityContext := other.SecurityContext
	newTemplate.SecurityContext = mySecurityContext
	if theirSecurityContext != nil {
		newTemplate.SecurityContext = theirSecurityContext
	}

	myArgs := mx.Args
	theirArgs := other.Args
	newTemplate.Args = myArgs
	if theirArgs != nil {
		newTemplate.Args = theirArgs
	}


	sidecarMapping := make(map[string]*SidecarSpec, 0)
	mySidecars := mx.Sidecar
	theirSidecars := other.Sidecar
	if len(mySidecars) > 0 {
		for _, sidecar := range mySidecars {
			sidecarMapping[sidecar.Name] = sidecar
		}
	}
	if len(theirSidecars) > 0 {
		for _, sidecar := range theirSidecars {
			sidecarMapping[sidecar.Name] = sidecar
		}
	}
	sidecars := make([]*SidecarSpec, 0)
	for _, v := range sidecarMapping {
		sidecars = append(sidecars, v)
	}
	newTemplate.Sidecar = sidecars

	return &newTemplate
}

func (mx *MixinTemplate) ReduceMerge(templates []*MixinTemplate) *MixinTemplate {
	var candidate *MixinTemplate
	for _, template := range templates {
		if candidate == nil {
			candidate = template
		} else {
			candidate = candidate.Merge(template)
		}
	}
	if candidate == nil {
		return nil
	}
	return candidate.Merge(mx)
}

func ReduceTemplates(templates []*MixinTemplate) *MixinTemplate {
	var candidate *MixinTemplate
	for _, template := range templates {
		if candidate == nil {
			candidate = template
		} else {
			candidate = candidate.Merge(template)
		}
	}
	return candidate
}
