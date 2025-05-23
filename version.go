package manifest

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// ZiplineeVersion is the object that determines how version numbers are generated
type ZiplineeVersion struct {
	SemVer *ZiplineeSemverVersion `yaml:"semver,omitempty" json:",omitempty"`
	Custom *ZiplineeCustomVersion `yaml:"custom,omitempty" json:",omitempty"`
}

// UnmarshalYAML customizes unmarshalling an ZiplineeVersion
func (version *ZiplineeVersion) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		SemVer *ZiplineeSemverVersion `yaml:"semver"`
		Custom *ZiplineeCustomVersion `yaml:"custom"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	version.SemVer = aux.SemVer
	version.Custom = aux.Custom

	// set default property values
	version.SetDefaults()

	return nil
}

// SetDefaults sets default values for properties of ZiplineeVersion if not defined
func (version *ZiplineeVersion) SetDefaults() {
	if version.Custom == nil && version.SemVer == nil {
		version.SemVer = &ZiplineeSemverVersion{}
	}

	// if version is semver set defaults
	if version.SemVer != nil {
		if version.SemVer.Patch == "" {
			version.SemVer.Patch = "{{auto}}"
		}
		if version.SemVer.LabelTemplate == "" {
			version.SemVer.LabelTemplate = "{{branch}}"
		}
		if len(version.SemVer.ReleaseBranch.Values) == 0 {
			version.SemVer.ReleaseBranch.Values = []string{"master", "main"}
		}
	}

	// if version is custom set defaults
	if version.Custom != nil {
		if version.Custom.LabelTemplate == "" {
			version.Custom.LabelTemplate = "{{revision}}"
		}
	}
}

// Version returns the version number as a string
func (version *ZiplineeVersion) Version(params ZiplineeVersionParams) string {
	if version.Custom != nil {
		return version.Custom.Version(params)
	}
	if version.SemVer != nil {
		return version.SemVer.Version(params)
	}
	return ""
}

// ZiplineeCustomVersion represents a custom version using a template
type ZiplineeCustomVersion struct {
	LabelTemplate string `yaml:"labelTemplate"`
}

// Version returns the version number as a string
func (v *ZiplineeCustomVersion) Version(params ZiplineeVersionParams) string {
	return parseTemplate(v.LabelTemplate, params.GetFuncMap())
}

// ZiplineeSemverVersion represents semantic versioning (http://semver.org/)
type ZiplineeSemverVersion struct {
	Major         int                 `yaml:"major"`
	Minor         int                 `yaml:"minor"`
	Patch         string              `yaml:"patch"`
	LabelTemplate string              `yaml:"labelTemplate"`
	ReleaseBranch StringOrStringArray `yaml:"releaseBranch"`
}

// Version returns the version number as a string
func (v *ZiplineeSemverVersion) Version(params ZiplineeVersionParams) string {

	patchWithLabel := v.GetPatchWithLabel(params)

	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, patchWithLabel)
}

// GetPatchWithLabel returns the formatted patch and label
func (v *ZiplineeSemverVersion) GetPatchWithLabel(params ZiplineeVersionParams) string {

	patch := v.GetPatch(params)
	label := v.GetLabel(params)

	if v.ReleaseBranch.Contains(params.Branch) || label == "" {
		return patch
	}

	return fmt.Sprintf("%v-%v", patch, label)
}

// GetPatch returns the formatted patch
func (v *ZiplineeSemverVersion) GetPatch(params ZiplineeVersionParams) string {

	return parseTemplate(v.Patch, params.GetFuncMap())
}

// GetLabel returns the formatted label
func (v *ZiplineeSemverVersion) GetLabel(params ZiplineeVersionParams) string {

	label := parseTemplate(v.LabelTemplate, params.GetFuncMap())

	if startsWithNumber, _ := regexp.Match(`^[0-9]`, []byte(label)); startsWithNumber {

		// get first placeholder from label template to use as prefix
		re := regexp.MustCompile(`{{([^}]+)}}`)
		match := re.FindStringSubmatch(v.LabelTemplate)

		prefix := "label-"
		if len(match) > 1 {
			prefix = match[1] + "-"
		}

		return v.tidyLabel(prefix + label)
	}

	return v.tidyLabel(label)
}

func (v *ZiplineeSemverVersion) tidyLabel(label string) string {
	// in order for the label to be used as a dns label (part between dots) it should only use
	// lowercase letters, digits and hyphens and have a max length of 63 characters;
	// also it should start with a letter and not end in a hyphen

	// ensure the label is lowercase
	label = strings.ToLower(label)

	// replace all invalid characters with a hyphen
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	label = reg.ReplaceAllString(label, "-")

	// replace double hyphens with a single one
	label = strings.Replace(label, "--", "-", -1)

	// trim hyphens from start and end
	label = strings.Trim(label, "-")

	// ensure it starts with a letter, not a digit or hyphen
	reg = regexp.MustCompile(`^[0-9-]+`)
	label = reg.ReplaceAllString(label, "")

	if len(label) > 63 {
		label = label[:63]
	}

	return label
}

// ZiplineeVersionParams contains parameters used to generate a version number
type ZiplineeVersionParams struct {
	AutoIncrement int
	Branch        string
	Revision      string
}

// GetFuncMap returns ZiplineeVersionParams as a function map for use in templating
func (p *ZiplineeVersionParams) GetFuncMap() template.FuncMap {

	return template.FuncMap{
		"auto":     func() string { return fmt.Sprint(p.AutoIncrement) },
		"branch":   func() string { return p.Branch },
		"revision": func() string { return p.Revision },
	}
}

func parseTemplate(templateText string, funcMap template.FuncMap) string {
	tmpl, err := template.New("version").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return err.Error()
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return err.Error()
	}

	return buf.String()
}
