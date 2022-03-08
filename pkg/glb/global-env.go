package glb

import (
	"github.com/skhatri/shores/pkg/applog"
	"github.com/skhatri/shores/pkg/environment"
	"github.com/skhatri/shores/pkg/functions"
	"sort"
)

//TODO substitute using subst map in go template
func LoadVarsWithSubstitution(files []string, subst map[string]string) map[string]map[string]string {
	result := loadEnvData(files)
	data := make(map[string]map[string]string, 0)
	for k, v := range result {
		data[k] = v
	}
	return data
}

func LoadVars(files []string) map[string]string {
	keys := make([]string, 0)
	result := loadEnvData(files)
	for k, _ := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	data := make(map[string]string, 0)
	for _, k := range keys {
		namedData := result[k]
		for attrib, value := range namedData {
			old, exists := data[attrib]
			if exists {
				applog.Tag("load-vars").WithAttribute("key", attrib).
					WithAttribute("new_value", value).WithAttribute("prev_value", old).
					Error("data for attribute %s already exists", attrib)
			}
			data[attrib] = value
		}
	}
	data["REGION"] = environment.Region()
	data["ENV_NAME"] = environment.EnvName()
	data["CLUSTER"] = environment.Cluster()
	return data
}

func loadEnvData(files []string) map[string]map[string]string {
	errors := make([]string, 0)
	variables := make(map[string]map[string]string, 0)
	for _, file := range files {
		envData := Environment{}
		err := functions.UnmarshalFile(file, &envData)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if envData.Kind != "Environment" {
			continue
		}
		data := make(map[string]string, 0)
		if matchBySelector(envData) {
			for _, kv := range envData.Spec.Data {
				data[kv.Name] = kv.Value
			}
		}
		variables[envData.Metadata.Name] = data
	}
	if len(errors) > 0 {
		applog.Tag("load-vars").Error("errors while loading environment data: %s", errors)
	}
	return variables
}

func matchBySelector(envData Environment) bool {
	selector := envData.Spec.Selector
	includedByEnv := true
	includedByLoc := true
	include := true
	if len(selector) != 0 {
		targetEnv, ok := selector["ENV_NAME"]
		if ok {
			if targetEnv != environment.EnvName() {
				includedByEnv = false
			}
		}
		location, ok := selector["LOCATION"]
		if ok {
			if location != environment.Region() {
				includedByLoc = false
			}
		}
		include = includedByEnv && includedByLoc
	}
	return include
}
