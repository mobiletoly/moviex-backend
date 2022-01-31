{{- define "moviex.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "moviex.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "moviex.servicename" -}}
{{- $base := default .root.Chart.Name .root.Values.nameOverride -}}
{{- printf "%s-%s" $base .microService.service.name -}}
{{- end -}}
