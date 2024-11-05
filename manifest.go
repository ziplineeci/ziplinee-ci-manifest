package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"

	yaml "gopkg.in/yaml.v2"
)

// ZiplineeManifest is the object that the .ziplinee.yaml deserializes to
type ZiplineeManifest struct {
	Archived         bool                       `yaml:"archived,omitempty"`
	Builder          ZiplineeBuilder            `yaml:"builder,omitempty"`
	Labels           map[string]string          `yaml:"labels,omitempty"`
	Version          ZiplineeVersion            `yaml:"version,omitempty"`
	GlobalEnvVars    map[string]string          `yaml:"env,omitempty"`
	Triggers         []*ZiplineeTrigger         `yaml:"triggers,omitempty"`
	Stages           []*ZiplineeStage           `yaml:"-"`
	Releases         []*ZiplineeRelease         `yaml:"-"`
	ReleaseTemplates []*ZiplineeReleaseTemplate `yaml:"-"`
	Bots             []*ZiplineeBot             `yaml:"-"`
}

// UnmarshalYAML customizes unmarshalling an ZiplineeManifest
func (c *ZiplineeManifest) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Archived            bool               `yaml:"archived"`
		Builder             ZiplineeBuilder    `yaml:"builder"`
		Labels              map[string]string  `yaml:"labels"`
		Version             ZiplineeVersion    `yaml:"version"`
		GlobalEnvVars       map[string]string  `yaml:"env"`
		DeprecatedPipelines yaml.MapSlice      `yaml:"pipelines"`
		Triggers            []*ZiplineeTrigger `yaml:"triggers"`
		Stages              yaml.MapSlice      `yaml:"stages"`
		Releases            yaml.MapSlice      `yaml:"releases"`
		ReleaseTemplates    yaml.MapSlice      `yaml:"releaseTemplates"`
		Bots                yaml.MapSlice      `yaml:"bots"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	c.Archived = aux.Archived
	c.Builder = aux.Builder
	c.Version = aux.Version
	c.Labels = aux.Labels
	c.GlobalEnvVars = aux.GlobalEnvVars
	c.Triggers = aux.Triggers

	// provide backwards compatibility for the deprecated pipelines section now renamed to stages
	if len(aux.Stages) == 0 && len(aux.DeprecatedPipelines) > 0 {
		aux.Stages = aux.DeprecatedPipelines
	}

	for _, mi := range aux.Stages {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var stage *ZiplineeStage
		if err := yaml.Unmarshal(bytes, &stage); err != nil {
			return err
		}
		if stage == nil {
			stage = &ZiplineeStage{}
		}

		// set the stage name, overwriting the name property if set on the stage explicitly
		stage.Name = mi.Key.(string)

		c.Stages = append(c.Stages, stage)
	}

	releaseTemplates := map[string]*ZiplineeReleaseTemplate{}

	for _, mi := range aux.ReleaseTemplates {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var releaseTemplate *ZiplineeReleaseTemplate
		if err := yaml.Unmarshal(bytes, &releaseTemplate); err != nil {
			return err
		}
		if releaseTemplate == nil {
			releaseTemplate = &ZiplineeReleaseTemplate{}
		}

		if releaseTemplate.Name == "" {
			releaseTemplate.Name = mi.Key.(string)
		}
		c.ReleaseTemplates = append(c.ReleaseTemplates, releaseTemplate)

		releaseTemplates[releaseTemplate.Name] = releaseTemplate
	}

	for _, mi := range aux.Releases {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var release *ZiplineeRelease
		if err := yaml.Unmarshal(bytes, &release); err != nil {
			return err
		}
		if release == nil {
			release = &ZiplineeRelease{}
		}

		if release.Name == "" {
			release.Name = mi.Key.(string)
		}

		release.InitFromTemplate(releaseTemplates)

		c.Releases = append(c.Releases, release)
	}

	for _, mi := range aux.Bots {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var bot *ZiplineeBot
		if err := yaml.Unmarshal(bytes, &bot); err != nil {
			return err
		}
		if bot == nil {
			bot = &ZiplineeBot{}
		}

		bot.Name = mi.Key.(string)
		c.Bots = append(c.Bots, bot)
	}

	return nil
}

// MarshalYAML customizes marshalling an ZiplineeManifest
func (c ZiplineeManifest) MarshalYAML() (out interface{}, err error) {
	var aux struct {
		Archived         bool               `yaml:"archived,omitempty"`
		Builder          ZiplineeBuilder    `yaml:"builder,omitempty"`
		Labels           map[string]string  `yaml:"labels,omitempty"`
		Version          ZiplineeVersion    `yaml:"version,omitempty"`
		GlobalEnvVars    map[string]string  `yaml:"env,omitempty"`
		Triggers         []*ZiplineeTrigger `yaml:"triggers,omitempty"`
		Stages           yaml.MapSlice      `yaml:"stages,omitempty"`
		Releases         yaml.MapSlice      `yaml:"releases,omitempty"`
		ReleaseTemplates yaml.MapSlice      `yaml:"releaseTemplates,omitempty"`
		Bots             yaml.MapSlice      `yaml:"bots,omitempty"`
	}

	aux.Archived = c.Archived
	aux.Builder = c.Builder
	aux.Labels = c.Labels
	aux.Version = c.Version
	aux.GlobalEnvVars = c.GlobalEnvVars
	aux.Triggers = c.Triggers

	for _, stage := range c.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}
	for _, release := range c.Releases {
		aux.Releases = append(aux.Releases, yaml.MapItem{
			Key:   release.Name,
			Value: release,
		})
	}
	for _, releaseTemplate := range c.ReleaseTemplates {
		aux.ReleaseTemplates = append(aux.ReleaseTemplates, yaml.MapItem{
			Key:   releaseTemplate.Name,
			Value: releaseTemplate,
		})
	}
	for _, bot := range c.Bots {
		aux.Bots = append(aux.Bots, yaml.MapItem{
			Key:   bot.Name,
			Value: bot,
		})
	}

	return aux, err
}

// SetDefaults sets default values for properties of ZiplineeManifest if not defined
func (c *ZiplineeManifest) SetDefaults(preferences ZiplineeManifestPreferences) {
	c.Builder.SetDefaults(preferences)
	c.Version.SetDefaults()

	for _, t := range c.Triggers {
		t.SetDefaults(preferences, TriggerTypeBuild, "")
	}
	for _, s := range c.Stages {
		s.SetDefaults(c.Builder)
	}

	for _, r := range c.Releases {
		if r.CloneRepository == nil {
			falseValue := false
			r.CloneRepository = &falseValue
		}

		if r.Builder == nil {
			r.Builder = &c.Builder
		} else {
			r.Builder.SetDefaults(preferences)
		}
		for _, t := range r.Triggers {
			t.SetDefaults(preferences, TriggerTypeRelease, r.Name)
		}
		for _, s := range r.Stages {
			s.SetDefaults(*r.Builder)
		}
	}

	for _, b := range c.Bots {
		if b.CloneRepository == nil {
			falseValue := false
			b.CloneRepository = &falseValue
		}

		if b.Builder == nil {
			b.Builder = &c.Builder
		} else {
			b.Builder.SetDefaults(preferences)
		}
		for _, t := range b.Triggers {
			t.SetDefaults(preferences, TriggerTypeBot, b.Name)
		}
		for _, s := range b.Stages {
			s.SetDefaults(*b.Builder)
		}
	}
}

// Validate checks if the manifest is valid
func (c *ZiplineeManifest) Validate(preferences ZiplineeManifestPreferences) (err error) {

	err = c.Builder.validate(preferences)
	if err != nil {
		return
	}

	// loop labels and check if they meet the label regexes
	for key, value := range c.Labels {
		if pattern, ok := preferences.LabelRegexes[key]; ok {
			pattern = fmt.Sprintf("^%v$", strings.TrimSpace(pattern))

			match, err := regexp.MatchString(pattern, value)
			if err != nil {
				return err
			}

			if !match {
				return fmt.Errorf("Label %v does not match regex %v", key, pattern)
			}
		}
	}

	if len(c.Stages) == 0 {
		return fmt.Errorf("The manifest should define 1 or more stages")
	}
	for _, s := range c.Stages {
		err = s.Validate()
		if err != nil {
			return
		}
	}

	for _, t := range c.Triggers {
		err = t.Validate("build", "")
		if err != nil {
			return
		}
	}

	for _, r := range c.Releases {
		if r.Builder != nil {
			err = r.Builder.validate(preferences)
			if err != nil {
				return
			}
		}

		for _, t := range r.Triggers {
			err = t.Validate(TriggerTypeRelease, r.Name)
			if err != nil {
				return
			}
		}

		for _, s := range r.Stages {
			err = s.Validate()
			if err != nil {
				return
			}
		}
	}

	for _, b := range c.Bots {
		if b.Builder != nil {
			err = b.Builder.validate(preferences)
			if err != nil {
				return
			}
		}

		for _, t := range b.Triggers {
			err = t.Validate(TriggerTypeBot, b.Name)
			if err != nil {
				return
			}
		}

		for _, s := range b.Stages {
			err = s.Validate()
			if err != nil {
				return
			}
		}
	}

	return nil
}

// GetAllTriggers returns both build and release triggers as one list
func (c *ZiplineeManifest) GetAllTriggers(repoSource, repoOwner, repoName string) []ZiplineeTrigger {
	// collect both build and release triggers
	triggers := make([]ZiplineeTrigger, 0)

	pipelineName := fmt.Sprintf("%v/%v/%v", repoSource, repoOwner, repoName)

	// add all build triggers
	for _, t := range c.Triggers {
		if t != nil {
			t.ReplaceSelf(pipelineName)
			triggers = append(triggers, *t)
		}
	}

	// add all release triggers
	for _, r := range c.Releases {
		for _, t := range r.Triggers {
			if t != nil {
				t.ReplaceSelf(pipelineName)
				triggers = append(triggers, *t)
			}
		}
	}

	// add all bots triggers
	for _, b := range c.Bots {
		for _, t := range b.Triggers {
			if t != nil {
				t.ReplaceSelf(pipelineName)
				triggers = append(triggers, *t)
			}
		}
	}

	return triggers
}

// DeepCopy provides a copy of all nested pointers
func (c *ZiplineeManifest) DeepCopy() (target ZiplineeManifest) {

	copier.CopyWithOption(&target, c, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return
}

// Exists checks whether the .ziplinee.yaml exists
func Exists(manifestPath string) bool {

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// does not exist
		return false
	}

	// does exist
	return true
}

// ReadManifestFromFile reads the .ziplinee.yaml into an ZiplineeManifest object
func ReadManifestFromFile(preferences *ZiplineeManifestPreferences, manifestPath string, validate bool) (manifest ZiplineeManifest, err error) {

	log.Debug().Msgf("Reading %v file...", manifestPath)

	// default preferences if not passed
	if preferences == nil {
		preferences = GetDefaultManifestPreferences()
	}

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return manifest, err
	}

	// unmarshal strict, so non-defined properties or incorrect nesting will fail
	if err := yaml.UnmarshalStrict(data, &manifest); err != nil {
		return manifest, err
	}

	// set defaults
	manifest.SetDefaults(*preferences)

	if validate {
		// check if manifest is valid
		err = manifest.Validate(*preferences)
		if err != nil {
			return manifest, err
		}
	}

	log.Debug().Msgf("Finished reading %v file successfully", manifestPath)

	return
}

// ReadManifest reads the string representation of .ziplinee.yaml into an ZiplineeManifest object
func ReadManifest(preferences *ZiplineeManifestPreferences, manifestString string, validate bool) (manifest ZiplineeManifest, err error) {

	// default preferences if not passed
	if preferences == nil {
		preferences = GetDefaultManifestPreferences()
	}

	// unmarshal strict, so non-defined properties or incorrect nesting will fail
	if err := yaml.UnmarshalStrict([]byte(manifestString), &manifest); err != nil {
		return manifest, err
	}

	// set defaults
	manifest.SetDefaults(*preferences)

	if validate {
		// check if manifest is valid
		err = manifest.Validate(*preferences)
		if err != nil {
			return manifest, err
		}
	}

	return
}

// GetDefaultManifestPreferences returns default preferences if not configured at the server
func GetDefaultManifestPreferences() (preferences *ZiplineeManifestPreferences) {

	preferences = &ZiplineeManifestPreferences{}
	preferences.SetDefaults()

	return
}
