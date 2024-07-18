// Code generated by internal/generate/tagstests/main.go; DO NOT EDIT.

{{ define "Init" }}
	ctx := acctest.Context(t)
	dataSourceName := "data.{{ .TypeName}}.test"{{ if .Generator }}
	rName := {{ .Generator }}
{{- end }}
{{ range .InitCodeBlocks -}}
{{ .Code }}
{{- end }}
{{ end }}

{{ define "TestCaseSetup" -}}
{{ template "TestCaseSetupNoProviders" . -}}
{{ if not .AlternateRegionProvider }}
	ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
{{- end -}}
{{- end }}

{{ define "TestCaseSetupNoProviders" -}}
	PreCheck:   func() { acctest.PreCheck(ctx, t){{ if .PreCheck }}; testAccPreCheck(ctx, t){{ end }} },
	ErrorCheck: acctest.ErrorCheck(t, names.{{ .PackageProviderNameUpper }}ServiceID),
{{- end }}

{{ define "TagsKnownValueForNull" -}}
{{ if eq .Implementation "framework" -}}
knownvalue.Null()
{{- else -}}
	{{ if .TagsUpdateForceNew -}}
	knownvalue.Null()
	{{- else -}}
	knownvalue.MapExact(map[string]knownvalue.Check{})
	{{- end }}
{{- end }}
{{- end }}

{{ define "TagsAllPlanCheckForNull" -}}
{{ if .TagsUpdateForceNew -}}
plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New(names.AttrTagsAll)),
{{- else -}}
plancheck.ExpectKnownValue(resourceName, tfjsonpath.New(names.AttrTagsAll), knownvalue.MapExact(map[string]knownvalue.Check{})),
{{- end }}
{{- end }}

{{ define "ImportBody" }}
	ResourceName: resourceName,
	ImportState:  true,
{{ if gt (len .ImportStateID) 0 -}}
	ImportStateId: {{ .ImportStateID }},
{{ end -}}
{{ if gt (len .ImportStateIDFunc) 0 -}}
	ImportStateIdFunc: {{ .ImportStateIDFunc }}(resourceName),
{{ end -}}
	ImportStateVerify: true,
{{ if gt (len .ImportIgnore) 0 -}}
	ImportStateVerifyIgnore: []string{
	{{ range $i, $v := .ImportIgnore }}{{ $v }},{{ end }}
	},
{{- end }}
{{ end }}

{{ define "testname" -}}
{{ if .Serialize }}testAcc{{ else }}TestAcc{{ end }}{{ .ResourceProviderNameUpper }}{{ .Name }}DataSource
{{- end }}

{{ define "ExistsCheck" }}
	testAccCheck{{ .Name }}Exists(ctx, {{ if .ExistsTakesT }}t,{{ end }} resourceName{{ if .ExistsTypeName}}, &v{{ end }}),
{{ end }}

{{ define "AdditionalTfVars" -}}
	{{ range $name, $value := .AdditionalTfVars -}}
	{{ $name }}: config.StringVariable({{ $value }}),
	{{ end -}}
	{{ if .AlternateRegionProvider -}}
	"alt_region": config.StringVariable(acctest.AlternateRegion()),
	{{ end }}
{{ end }}

package {{ .ProviderPackage }}_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
	{{ range .GoImports -}}
	{{ if .Alias }}{{ .Alias }} {{ end }}"{{ .Path }}"
	{{ end }}
)

{{ if .Serialize }}
func {{ template "testname" . }}_tagsSerial(t *testing.T) {
	t.Helper()

	testCases := map[string]func(t *testing.T){
		acctest.CtBasic: {{ template "testname" . }}_tags,
	}

	acctest.RunSerialTests1Level(t, testCases, {{ if .SerializeDelay }}serializeDelay{{ else }}0{{ end }})
}
{{ end }}

func {{ template "testname" . }}_tags(t *testing.T) {
	{{- template "Init" . }}

	resource.{{ if .Serialize }}Test{{ else }}ParallelTest{{ end }}(t, resource.TestCase{
		{{ template "TestCaseSetup" . }}
		Steps: []resource.TestStep{
			{
				{{ if .AlternateRegionProvider -}}
				ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
				{{ end -}}
				ConfigDirectory: config.StaticDirectory("testdata/{{ .Name }}/data.tags/"),
				ConfigVariables: config.Variables{ {{ if .Generator }}
					acctest.CtRName: config.StringVariable(rName),{{ end }}
					acctest.CtResourceTags: config.MapVariable(map[string]config.Variable{
						acctest.CtKey1: config.StringVariable(acctest.CtValue1),
					}),
					{{ template "AdditionalTfVars" . }}
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New(names.AttrTags), knownvalue.MapExact(map[string]knownvalue.Check{
						acctest.CtKey1: knownvalue.StringExact(acctest.CtValue1),
					})),
				},
			},
		},
	})
}
