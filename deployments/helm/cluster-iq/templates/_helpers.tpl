{{/*
Expand the name of the chart.
*/}}
{{- define "cluster-iq.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cluster-iq.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cluster-iq.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cluster-iq.labels" -}}
helm.sh/chart: {{ include "cluster-iq.chart" . }}
{{ include "cluster-iq.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Component labels
*/}}
{{- define "cluster-iq.componentLabels" -}}
app.kubernetes.io/component: {{ . }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cluster-iq.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cluster-iq.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the API service account
*/}}
{{- define "cluster-iq.apiServiceAccountName" -}}
{{- if .Values.api.serviceAccount.create }}
{{- default "api" .Values.api.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.api.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the Console service account
*/}}
{{- define "cluster-iq.consoleServiceAccountName" -}}
{{- if .Values.console.serviceAccount.create }}
{{- default "console" .Values.console.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.console.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the Scanner service account
*/}}
{{- define "cluster-iq.scannerServiceAccountName" -}}
{{- if .Values.scanner.serviceAccount.create }}
{{- default "scanner" .Values.scanner.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.scanner.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the Agent service account
*/}}
{{- define "cluster-iq.agentServiceAccountName" -}}
{{- if .Values.agent.serviceAccount.create }}
{{- default "agent" .Values.agent.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.agent.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the Database service account
*/}}
{{- define "cluster-iq.databaseServiceAccountName" -}}
{{- if .Values.database.serviceAccount.create }}
{{- default "database" .Values.database.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.database.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "cluster-iq.backendUrl" -}}
{{- printf "http://api.%s.svc.cluster.local:%v" .Release.Namespace .Values.api.service.port -}}
{{- end }}



