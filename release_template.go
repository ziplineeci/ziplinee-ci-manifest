package manifest

import (
	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v2"
)

// ZiplineeReleaseTemplate represents a template for a release target
type ZiplineeReleaseTemplate struct {
	Name            string                   `yaml:"-"`
	Builder         *ZiplineeBuilder         `yaml:"builder,omitempty"`
	CloneRepository *bool                    `yaml:"clone,omitempty" json:",omitempty"`
	Actions         []*ZiplineeReleaseAction `yaml:"actions,omitempty" json:",omitempty"`
	Triggers        []*ZiplineeTrigger       `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*ZiplineeStage         `yaml:"-"`
}

// UnmarshalYAML customizes unmarshalling an ZiplineeRelease
func (releaseTemplate *ZiplineeReleaseTemplate) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string                   `yaml:"name"`
		Builder         *ZiplineeBuilder         `yaml:"builder"`
		CloneRepository *bool                    `yaml:"clone"`
		Actions         []*ZiplineeReleaseAction `yaml:"actions"`
		Triggers        []*ZiplineeTrigger       `yaml:"triggers"`
		Stages          yaml.MapSlice            `yaml:"stages"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	releaseTemplate.Name = aux.Name
	releaseTemplate.Builder = aux.Builder
	releaseTemplate.CloneRepository = aux.CloneRepository
	releaseTemplate.Actions = aux.Actions
	releaseTemplate.Triggers = aux.Triggers

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

		releaseTemplate.Stages = append(releaseTemplate.Stages, stage)
	}

	return nil
}

// MarshalYAML customizes marshalling an ZiplineeManifest
func (releaseTemplate ZiplineeReleaseTemplate) MarshalYAML() (out interface{}, err error) {

	var aux struct {
		Name            string                   `yaml:"-"`
		Builder         *ZiplineeBuilder         `yaml:"builder,omitempty"`
		CloneRepository *bool                    `yaml:"clone,omitempty"`
		Actions         []*ZiplineeReleaseAction `yaml:"actions,omitempty"`
		Triggers        []*ZiplineeTrigger       `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice            `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = releaseTemplate.Builder
	aux.CloneRepository = releaseTemplate.CloneRepository
	aux.Actions = releaseTemplate.Actions
	aux.Triggers = releaseTemplate.Triggers

	for _, stage := range releaseTemplate.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}

// DeepCopy provides a copy of all nested pointers
func (releaseTemplate ZiplineeReleaseTemplate) DeepCopy() (target ZiplineeReleaseTemplate) {

	copier.CopyWithOption(&target, releaseTemplate, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return
}
