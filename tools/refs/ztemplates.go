// Code generated by assets. DO NOT EDIT.

package main

import "fmt"

func loadtemplate(name string) ([]byte, error) {
	switch name {
	case "markdown.tmpl":
		return []byte(`# References
{{ $references := .References -}}
{{ range $_, $sec :=.Sections -}}
## {{ $sec.Name }}
{{ range $_, $ref := $references -}}
{{ if and (not $ref.Queued) (eq $ref.Section $sec.ID) -}}
* [{{ if $ref.Highlight }}**{{ end }}{{ $ref.Title }}{{ if $ref.Highlight }}**{{ end }}]({{ $ref.URL }})
    {{- if $ref.Supplements }} (
        {{- range $i, $supp := $ref.Supplements -}}
        {{- if ne $i 0 }}, {{ end -}}
        [{{ $supp.Type }}]({{ $supp.URL }})
        {{- end -}}
    ){{ end -}}
    {{- if $ref.Author }} {{ $ref.Author }}.{{- end}}
    {{- if $ref.Note }} _Note:_ {{ $ref.Note }}{{- end}}
{{ end -}}
{{- end -}}
{{- end -}}
`), nil

	default:
		return nil, fmt.Errorf("unknown asset %s", name)
	}
}
