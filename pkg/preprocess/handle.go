package preprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skhatri/shores/pkg/applog"
	"github.com/skhatri/shores/pkg/environment"
	"github.com/skhatri/shores/pkg/model"
	"strings"
)

func enrichAppSpecification(spec model.AppSpec, envLookupData map[string]map[string]string,
	resourceLookupData map[string]model.Resources) model.Deployable {

	targetInfo := createTargetInfo(spec)
	healthChecks := createChecks(spec.Service)
	services := createServices(spec.Service)
	envData := createEnv(spec.Env, envLookupData)
	serviceEnabled := len(services) > 0
	resources := createResources(spec.Resources, resourceLookupData)
	ingress := spec.Ingress
	return model.Deployable{
		Artifact: model.ArtifactInfo{
			Name:  spec.Name,
			Image: spec.Image,
		},
		Env:                envData,
		Checks:             healthChecks,
		Target:             targetInfo,
		Service:            services,
		ServiceAccountName: spec.ServiceAccount,
		ServiceEnabled:     serviceEnabled,
		Resources:          resources,
		Ingress:            ingress,
	}
}

func createResources(resources []string, data map[string]model.Resources) *model.Resources {
	var resourceRef = &model.Resources{}
	if len(resources) == 0 {
		resources = append(resources, "small")
	}

	for _, resource := range resources {
		resourceSpec := data[resource]
		if resourceSpec.Limits != nil {
			if resourceSpec.Limits.Cpu != nil {
				if resourceRef.Limits == nil {
					resourceRef.Limits = &model.ResourceValue{}
				}
				resourceRef.Limits.Cpu = resourceSpec.Limits.Cpu
			}
			if resourceSpec.Limits.Memory != nil {
				if resourceRef.Limits == nil {
					resourceRef.Limits = &model.ResourceValue{}
				}
				resourceRef.Limits.Memory = resourceSpec.Limits.Memory
			}

		}

		if resourceSpec.Requests != nil {
			if resourceSpec.Requests.Cpu != nil {
				if resourceRef.Requests == nil {
					resourceRef.Requests = &model.ResourceValue{}
				}
				resourceRef.Requests.Cpu = resourceSpec.Requests.Cpu
			}
			if resourceSpec.Requests.Memory != nil {
				if resourceRef.Requests == nil {
					resourceRef.Requests = &model.ResourceValue{}
				}
				resourceRef.Requests.Memory = resourceSpec.Requests.Memory
			}
		}
	}
	if resourceRef.Limits == nil && resourceRef.Requests == nil {
		return nil
	}
	return resourceRef
}

func createEnv(vars []model.Env, lookupData map[string]map[string]string) map[string]string {
	envData := make(map[string]string, 0)
	for _, v := range vars {
		if v.EnvSet != nil {
			envSetData, ok := lookupData[*v.EnvSet]
			if ok {
				for key, value := range envSetData {
					envData[key] = value
				}
			} else {
				applog.Tag("load-app-env").WithAttribute("env_set", *v.EnvSet).Error("Env set not found")
			}
		}
	}
	for _, v := range vars {
		if v.Name != nil && v.Value != nil {
			envData[*v.Name] = *v.Value
		}
	}
	return envData
}

func createChecks(service *model.ServiceSpec) *model.Healthcheck {
	if service == nil || service.HealthCheckUrl == nil {
		return nil
	}
	url := service.HealthCheckUrl
	healthCheck := model.Healthcheck{
		Port: "http",
		Path: *url,
	}
	return &healthCheck
}

func updateDeploymentArtifact(deploymentSpec *model.Deployable, releaseSpec model.ReleaseSpec) {
	deploymentSpec.Artifact = model.ArtifactInfo{
		Name:  releaseSpec.Name,
		Image: *releaseSpec.Image,
	}
	deploymentSpec.Namespace = releaseSpec.Namespace
}

func createServices(service *model.ServiceSpec) []model.ServiceInfo {
	services := make([]model.ServiceInfo, 0)
	if service == nil {
		return services
	}
	ports := make([]model.PortType, 0)
	for k, v := range service.Port {
		portTypeInstance := model.PortType{
			Port:       v,
			Name:       k,
			Protocol:   "TCP",
			TargetPort: k,
		}
		ports = append(ports, portTypeInstance)
	}
	info := model.ServiceInfo{
		Headless: false,
		Port:     ports,
	}
	services = append(services, info)
	if service.Headless != nil && *service.Headless {
		headless := model.ServiceInfo{
			Headless: true,
			Port:     ports,
		}
		services = append(services, headless)
	}
	return services
}

func mergeMixins(spec *model.AppSpec, mixinsData map[string]model.MixinTemplate) {
	mixins := make([]*model.MixinTemplate, 0)
	for _, mixinName := range spec.Mixins {
		mixinRef := mixinsData[mixinName]
		mixins = append(mixins, &mixinRef)
	}
	mixinTemplate := model.ReduceTemplates(mixins)

	if mixinTemplate != nil {
		if spec.Service == nil {
			spec.Service = mixinTemplate.Service
		}
		if spec.Workload == nil {
			spec.Workload = mixinTemplate.Workload
		}
		if spec.SecurityContext == nil {
			spec.SecurityContext = mixinTemplate.SecurityContext
		}
		if spec.Resources == nil {
			if mixinTemplate.Resources != nil {
				for _, resName := range mixinTemplate.Resources {
					spec.Resources = append(spec.Resources, *resName)
				}
			}
		}
		if spec.Secrets == nil {
			spec.Secrets = mixinTemplate.Secrets
		}
		if spec.Args == nil {
			spec.Args = mixinTemplate.Args
		}
	}
}

func createTargetInfo(spec model.AppSpec) model.TargetInfo {
	defaultScaling := "tools"
	if spec.Workload == nil {
		spec.Workload = &model.WorkloadSpec{
			Target:  "tools",
			Scaling: &defaultScaling,
		}
	}
	replica := 1
	if spec.Workload.Scaling == nil {
		spec.Workload.Scaling = &defaultScaling
	}
	if spec.Workload.Scaling != nil {
		replica = replicaForScalingGroup(*spec.Workload.Scaling)
	}
	targetInfo := model.TargetInfo{
		NodeSelector: map[string]string{
			"eks.amazonaws.com/nodegroup": spec.Workload.Target,
		},
		Replica: replica,
	}
	return targetInfo
}

func replicaForScalingGroup(group string) int {
	switch group {
	case "microservices":
		if environment.IsProd() {
			return 3
		}
		return 1
	}
	return 1
}
func ValidateAppSpec(spec model.AppSpec,
	envLookupData map[string]map[string]string,
	resourceLookupData map[string]model.Resources,
	mixinsData map[string]model.MixinTemplate,
	releaseSpec model.ReleaseSpec,
	task model.Task) (*model.Deployable, error) {

	mergeMixins(&spec, mixinsData)
	deploymentSpec := enrichAppSpecification(spec, envLookupData, resourceLookupData)
	updateDeploymentArtifact(&deploymentSpec, releaseSpec)
	updateLabelsAndAnnotations(&deploymentSpec, releaseSpec, task)
	updateSecurityContext(&deploymentSpec, spec)
	updateMount(&deploymentSpec, spec)
	updateArgs(&deploymentSpec, spec)
	return &deploymentSpec, nil
}

var (
	True  = true
	False = false
)

func updateSecurityContext(deployable *model.Deployable, appSpec model.AppSpec) {
	if appSpec.SecurityContext != nil {
		deployable.SecurityContext = appSpec.SecurityContext
	} else {
		deployable.SecurityContext = &model.SecurityContextSpec{
			RunAsUser:                "1000",
			AllowPrivilegeEscalation: &False,
			ReadOnlyRootFilesystem:   &True,
			RunAsNonRoot:             &True,
		}
	}
}

func updateLabelsAndAnnotations(deploymentSpec *model.Deployable, releaseSpec model.ReleaseSpec, task model.Task) {
	labels := make(map[string]string, 0)
	labels["helm.sh/chart"] = fmt.Sprintf("%s-%s", releaseSpec.Name, *releaseSpec.Version)
	labels["app.kubernetes.io/name"] = releaseSpec.Name
	labels["app.kubernetes.io/instance"] = releaseSpec.Name
	labels["app.kubernetes.io/managed-by"] = "Helm"
	labels["app.kubernetes.io/version"] = *releaseSpec.Version
	labels["app.kubernetes.io/release"] = releaseSpec.Name

	annotations := make(map[string]string, 0)
	str := bytes.Buffer{}
	json.NewEncoder(&str).Encode(releaseSpec)
	annotations["app.kubernetes.io/artifact-info"] = strings.ReplaceAll(str.String(), "\n", "")

	taskInfo := bytes.Buffer{}
	json.NewEncoder(&taskInfo).Encode(task)
	annotations["app.kubernetes.io/deployment-info"] = strings.ReplaceAll(taskInfo.String(), "\n", "")

	selectorLabels := map[string]string{
		"app.kubernetes.io/name":     releaseSpec.Name,
		"app.kubernetes.io/instance": releaseSpec.Name,
	}
	deploymentSpec.Metadata = model.Metadata{
		Labels:         labels,
		Annotations:    annotations,
		SelectorLabels: selectorLabels,
	}

}

func updateArgs(deployable *model.Deployable, spec model.AppSpec) {
	if spec.Args != nil {
		deployable.Args = spec.Args
	}
}

func updateMount(deployable *model.Deployable, appSpec model.AppSpec) {
	mounts := make([]model.MountSpec, 0)
	if len(appSpec.Mounts) == 0 {
		mounts = append(mounts, model.MountSpec{
			Name: "tmp-volume",
			Path: "/tmp",
			Type: "emptyDir",
		})
	}

	for _, mountValue := range appSpec.Mounts {
		parts := strings.Split(mountValue, ":")
		typeName := "emptyDir"
		if len(parts) > 1 {
			typeName = parts[1]
		}
		mountName := strings.ReplaceAll(parts[0], "/", "-")
		if mountName[0] == '-' {
			mountName = mountName[1:]
		}
		mounts = append(mounts, model.MountSpec{
			Name: mountName,
			Path: parts[0],
			Type: typeName,
		})
	}

	deployable.Mounts = mounts
}
