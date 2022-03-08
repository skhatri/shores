package templates

import (
	"errors"
	"fmt"
	"github.com/skhatri/shores/pkg/model"
	"strings"
	"text/template"
)

var ChartTemplate = `apiVersion: v1
description: A Helm chart for Kubernetes {{ .Artifact.Name }}
name: {{ .Artifact.Name }}
version: {{ .Artifact.Version }}

`

var ServiceAccountTemplate = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Artifact.Name | ToLower }}
  namespace: {{ .Namespace }}
  {{ if .Metadata.Annotations }}annotations:
{{ range $key, $value := .Metadata.Annotations }}{{ $key | indent 4 }}: '{{ $value }}'
{{ end }}{{ end }}
  {{ if .Metadata.Labels }}labels:
{{ range $key, $value := .Metadata.Labels }}{{ $key | indent 4 }}: '{{ $value }}'
{{ end }}{{ end }}

`

var ServiceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{ .Artifact.Name | ToLower }}
  namespace: {{ .Namespace }}
  {{ if .Metadata.Annotations }}annotations: 
{{ range $key, $value := .Metadata.Annotations }}{{ $key | indent 4 }}: '{{ $value }}'
{{ end }}{{ end }}
  {{ if .Metadata.Labels }}labels: 
{{ range $key, $value := .Metadata.Labels }}{{ $key | indent 4}}: '{{ $value }}'
{{ end }}{{ end }}
spec:
  type: ClusterIP
  {{ if .ServiceEnabled }}ports: 
{{ range $port := .Ports }}
    - name: {{ $port.Name }}
      port: {{ $port.Port }}
      targetPort: {{ $port.Name }}
      protocol: {{ $port.Protocol }}
{{ end }}{{ end }}
  {{ if .Metadata.SelectorLabels }}selector:
{{ range $key, $value := .Metadata.SelectorLabels }}{{ $key | indent 4 }}: {{ $value }}
{{ end }}{{ end }}

`

var DeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Artifact.Name | ToLower }}
  namespace: {{ .Namespace }}
  {{ if .Metadata.Annotations }}annotations:
{{ range $key, $value := .Metadata.Annotations }}{{ $key |indent 4}}: '{{ $value }}'
{{ end }}{{ end }}
  {{ if .Metadata.Labels }}labels:
{{ range $key, $value := .Metadata.Labels }}{{ $key |indent 4 }}: '{{ $value }}'
{{ end }}{{ end }}
spec:
  replicas: {{ .Target.Replica }}
  selector:
    {{ if .Metadata.SelectorLabels }}matchLabels:
{{ range $key, $value := .Metadata.SelectorLabels }}{{ $key | indent 6 }}: {{ $value }}
{{ end }}{{ end }}
  template:
    metadata:
      {{ if .Metadata.SelectorLabels }}labels:
{{ range $key, $value := .Metadata.SelectorLabels }}{{ $key | indent 8 }}: {{ $value }}
{{ end }}{{ end }}
    spec:
      serviceAccountName: {{ if .ServiceAccountName }}{{ .ServiceAccountName }}{{else}}{{ .Artifact.Name | ToLower }}{{end}}
      containers:
        - name: {{ .Artifact.Name }}
          image: {{ .Artifact.Image }}
          imagePullPolicy: IfNotPresent
{{ if .Args }}
          {{ if .Args.Entrypoint }}command:{{ range $entry := .Args.Entrypoint }} 
            - '{{$entry}}'{{ end }}{{ end }}
          {{ if .Args.Command }}args:{{ range $cmd := .Args.Command }}
            - '{{$cmd}}'{{ end }} {{ end }}
{{ end }}
          {{ if .ServiceEnabled }}ports:{{ range $port := .Ports }}
            - name: {{ $port.Name }}
              containerPort: {{ $port.Port }}
              protocol: {{ $port.Protocol }}
{{- end }}{{- end }}
          {{ if .Checks -}}livenessProbe:
            httpGet:
              path: {{ .Checks.Path }}
              port: {{ .Checks.Port }}
              initialDelaySeconds: 30
              timeoutSeconds: 100{{- end }}
          {{ if .Checks -}}readinessProbe:
            httpGet:
              path: {{ .Checks.Path }}
              port: {{ .Checks.Port }}
              initialDelaySeconds: 30
              timeoutSeconds: 100 {{- end }}
          {{if .Resources}}resources:
            {{ if .Resources.Requests }}requests:
              {{ if .Resources.Requests.Cpu }}cpu: "{{ .Resources.Requests.Cpu }}"{{end}}
              {{ if .Resources.Requests.Memory}}memory: "{{ .Resources.Requests.Memory }}"{{end}}
		    {{- end }}
            {{ if .Resources.Limits }}limits:
              {{ if .Resources.Limits.Cpu }}cpu: "{{ .Resources.Limits.Cpu }}"{{end}}
              {{ if .Resources.Limits.Memory}}memory: "{{ .Resources.Limits.Memory }}"{{end}}
            {{- end }}
		  {{- end }}
          {{ if .Env }}env:{{ range $key, $value := .Env }}
            - name: "{{ $key | ToUpper }}"
              value: "{{ $value }}"{{end}}
		  {{- end }}
          {{ if .Mounts }}volumeMounts:{{ range $mount := .Mounts }}
            - name: {{ $mount.Name }}
              mountPath: {{ $mount.Path }}
{{- end }}{{- end}}

          {{ if .SecurityContext}}securityContext:
            {{ if .SecurityContext.AllowPrivilegeEscalation }}allowPrivilegeEscalation: {{ .SecurityContext.AllowPrivilegeEscalation }}{{ end }}
            {{ if .SecurityContext.ReadOnlyRootFilesystem }}readOnlyRootFilesystem: {{ .SecurityContext.ReadOnlyRootFilesystem}}{{end}}
            {{ if .SecurityContext.RunAsNonRoot }}runAsNonRoot: {{ .SecurityContext.RunAsNonRoot }}{{ end }}
            {{ if .SecurityContext.RunAsUser }}runAsUser: {{ .SecurityContext.RunAsUser }}{{ end }}
          {{- end }}
      {{ if .SecurityContext}}securityContext:
        {{ if .SecurityContext.RunAsNonRoot }}runAsNonRoot: {{ .SecurityContext.RunAsNonRoot }}{{ end }}
        {{ if .SecurityContext.RunAsUser }}runAsUser: {{ .SecurityContext.RunAsUser }}{{ end }}
      {{- end }}
      {{ if .Target.NodeSelector }}nodeSelector:
{{ range $key, $value := .Target.NodeSelector }}{{ $key | indent 8 }}: {{ $value}}
{{end}}{{end}}
      affinity: { }
      tolerations: [ ]
      {{ if .Mounts }}volumes:{{ range $mount := .Mounts }}
        - name: {{ $mount.Name }}
          {{ if eq $mount.Type "emptyDir" }}emptyDir: { }{{end}}
{{- end }}{{- end}}
`

var JobTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
    version: {{ .Tag }}
  annotations:{{ if .Annotations }}
    {{ range $key, $value := .Annotations }}{{ $key }}: {{ $value }}
    {{ end }}{{ end }}
spec:
  {{ if .Replicas -}}completions: {{ .Replicas }}{{else}}1{{- end }}
  {{ if .Parallelism -}}parallelism: {{ .Parallelism }}{{- end }}
  {{ if .BackoffLimit -}}backoffLimit: {{ .BackoffLimit }}{{- end }}
  {{ if .ActiveDeadLine -}}activeDeadlineSeconds: {{ .ActiveDeadLine }}{{- end }}
  {{ if .TTLSecondsAfterFinished -}}ttlSecondsAfterFinished: {{ .TTLSecondsAfterFinished }}{{- end }}
  template:
    spec:
      serviceAccountName: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
      containers:
       - name: {{ .Name }}
         image: {{ .Name}}:{{ .Tag}}
         imagePullPolicy: IfNotPresent
         {{ if .Entrypoint }}command: [{{ range $entry := .Entrypoint }}'{{$entry}}', {{ end }}]{{ end }}
         {{ if .Command }}args: [{{ range $cmd := .Command }}'{{$cmd}}', {{ end }}]{{ end }}

         resources:
           limits:
             cpu: "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"   
           requests:
             cpu:  "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"
         {{ if .EnvVars }}env:{{ range $key, $value := .EnvVars }}
          - name: "{{ $key | ToUpper }}"
            value: "{{ $value }}"{{end}}
		 {{ end }}
      restartPolicy: {{ if .RestartPolicy -}}{{ .RestartPolicy }}{{ else }}Never{{end}} 
      affinity: {}
      nodeSelector: {}
      tolerations: []
`

//LoadTemplates parse static template to helm chart
func LoadTemplates(tName string, deployable *model.Deployable) (*template.Template, error) {
	appName := deployable.Artifact.Name
	switch tName {
	case "ChartTemplate":
		return getTemplate("Chart.yaml", ChartTemplate)
	case "DeploymentTemplate":
		return getTemplate(fmt.Sprintf("%s-deployment.yaml", appName), DeploymentTemplate)
	case "ServiceTemplate":
		return getTemplate(fmt.Sprintf("%s-service.yaml", appName), ServiceTemplate)
	case "ServiceAccountTemplate":
		return getTemplate(fmt.Sprintf("%s-serviceaccount.yaml", appName), ServiceAccountTemplate)
	case "JobTemplate":
		return getTemplate(fmt.Sprintf("%s-job.yaml", appName), JobTemplate)
	}
	return nil, nil
}

func getTemplate(name string, templateType string) (*template.Template, error) {
	indentFunc := func(suffix string) func(n int, s string) string {
		return func(n int, s string) string {
			return fmt.Sprintf("%s%s%s", strings.Repeat(" ", n), s, suffix)
		}
	}
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"indent":  indentFunc(""),
		"nindent": indentFunc("\n"),
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(templateType)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing %v ", err))
	}
	return tmpl, nil
}
