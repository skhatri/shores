package main

import (
	"fmt"
	"github.com/skhatri/shores/pkg/applog"
	"github.com/skhatri/shores/pkg/model"
	templates "github.com/skhatri/shores/pkg/template"
	"log"
	"os"
	"time"
)

func main() {
	release := "spec/user/release-set/release-1.yaml"
	now := time.Now()
	releaseId := now.Format("200601021504")
	task := model.Task{
		Action:    "deploy",
		Command:   fmt.Sprintf("deploy %s", release),
		ReleaseId: releaseId,
		User:      os.Getenv("USER"),
		Created:   now.Format(time.RFC3339),
		ChangeRef: "CRQ000019921",
		Output:    "../shores-helm/charts",
	}
	productSet, err := model.NewProductSetFromFile(release, "default")
	if err != nil {
		log.Fatalf("error processing product set file: %v", err)
	}
	dSummary, tErr := templates.Run(productSet, task)
	if tErr != nil {
		log.Fatalf("error running template: %v", tErr)
	}
	for _, item := range dSummary.Items {
		applog.Tag("summary").Info("generated for %s at %s", item.Name, item.Path)
	}
}
