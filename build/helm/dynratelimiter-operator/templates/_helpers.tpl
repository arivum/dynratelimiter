{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "dynratelimiter-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dynratelimiter-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "dynratelimiter-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "dynratelimiter-operator.labels" -}}
app.kubernetes.io/name: {{ include "dynratelimiter-operator.name" . }}
helm.sh/chart: {{ include "dynratelimiter-operator.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Webhook cert generation
*/}}
{{- define "dynratelimiter-operator.certs" -}}
{{- if not .Values.caCert -}}
{{- $ca := (genCA "enveloped-secrets-operator" 3650) -}}
{{- $svc1 := printf "%s" .Values.appname }}
{{- $svc2 := printf "%s.%s" .Values.appname .Values.namespace }}
{{- $svc3 := printf "%s.%s.svc" .Values.appname .Values.namespace }}
{{- $svc4 := printf "%s.%s.svc.cluster.local" .Values.appname .Values.namespace }}
{{- $cert := genSignedCert "enveloped-secrets-operator.arivum.de" (list) (list $svc1 $svc2 $svc3 $svc4) 3650 $ca -}}
{{- $_ := set .Values "caCert" ($ca.Cert | b64enc) -}}
{{- $_ := set .Values "tlsCert" ($cert.Cert | b64enc) -}}
{{- $_ := set .Values "tlsKey" ($cert.Key | b64enc) -}}
{{- end -}}
{{- end -}}