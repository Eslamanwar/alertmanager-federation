{{/*
Expand the name of the chart.
*/}}
{{- define "common.fullname" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "common.labels.standard" -}}
app.kubernetes.io/name: {{ include "common.fullname" . }}
{{- end -}}
