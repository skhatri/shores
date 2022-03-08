package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skhatri/shores/pkg/applog"
	"github.com/skhatri/shores/pkg/functions"
	"github.com/skhatri/shores/pkg/glb"
	"github.com/skhatri/shores/pkg/mixin"
	model "github.com/skhatri/shores/pkg/model"
	"github.com/skhatri/shores/pkg/preprocess"
	"github.com/skhatri/shores/pkg/resource"
	"os"
	"path/filepath"
	"strings"
)

func createDirSafely(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func Run(productSet *model.ProductSet, task model.Task, outputDir string) (*model.DeploymentSummary, error) {

	globalEnvData := glb.LoadVars(functions.ListFiles("spec/provider/globals", ".yaml"))
	envData := glb.LoadVarsWithSubstitution(functions.ListFiles("spec/provider/env-sets", ".yaml"), globalEnvData)

	resourcesData := resource.LoadResources(functions.ListFiles("spec/provider/resources", ".yaml"))
	mixinData := mixin.LoadMixins(functions.ListFiles("spec/provider/mixins", ".yaml"))

	itemSummary := model.DeploymentSummary{}
	items := make([]model.DeploymentItem, 0)
	for _, app := range productSet.Apps {
		appSpec := model.AppSpec{}
		uerr := functions.UnmarshalFile(fmt.Sprintf("spec/user/apps/%s.yaml", app.Name), &appSpec)
		if uerr != nil {
			return nil, uerr
		}
		applog.Tag("generator").WithAttribute("app_name", app.Name).Info("Generating app")
		deployable, err := preprocess.ValidateAppSpec(appSpec, envData, resourcesData, mixinData, *app, task)
		if err != nil {
			applog.Tag("deployer").WithAttribute("app_name", app.Name).Error("error %v", err)
			break
		}
		if applog.IsDebugEnabled() {
			b, e := json.Marshal(deployable)
			if e != nil {
				applog.Tag("marshaller").Error("error marshalling json", e)
			}
			fmt.Println(string(b))
		}

		appWorkDir := fmt.Sprintf("%s/%s/", outputDir, app.Name)
		cerr := createDirSafely(appWorkDir)
		if cerr != nil {
			return nil, cerr
		}
		requiredTemplates, kind := GetRequiredTemplates(deployable)
		for _, tName := range requiredTemplates {
			tmpl, err := LoadTemplates(tName, deployable)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("task: load-template, app: [%s], error: [%v]", app.Name, err))
			}
			file, er := os.Create(fmt.Sprintf("%s/%s", appWorkDir, tmpl.Name()))
			if er != nil {
				return nil, er
			}

			exErr := tmpl.Execute(file, &deployable)
			if exErr != nil {
				return nil, errors.New(fmt.Sprintf("task: execute, template: [%s], app: [%s], error: [%v]", tName, app.Name, err))
			}
		}
		items = append(items, model.DeploymentItem{
			Name: app.Name,
			Kind: kind,
			Path: appWorkDir,
		})
	}
	itemSummary = model.DeploymentSummary{
		Namespace: *productSet.Namespace,
		Items:     items,
	}
	return &itemSummary, nil
}

func GetRequiredTemplates(deployable *model.Deployable) ([]string, string) {
	kind := ""
	requiredTemplates := make([]string, 0)
	requiredTemplates = append(requiredTemplates, "ServiceAccountTemplate")
	if len(deployable.Kind) == 0 || deployable.Kind == "Deployment" {
		requiredTemplates = append(requiredTemplates, "DeploymentTemplate")
		kind = "deployment"
	} else if strings.EqualFold(deployable.Kind, "Job") {
		requiredTemplates = append(requiredTemplates, "JobTemplate")
		kind = "job"
	}

	if deployable.ServiceEnabled {
		requiredTemplates = append(requiredTemplates, "ServiceTemplate")
	}
	return requiredTemplates, kind
}
