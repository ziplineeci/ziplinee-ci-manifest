package manifest

import (
	yaml "gopkg.in/yaml.v2"
)

// ZiplineeBot allows to respond to any event coming from one of the integrations
type ZiplineeBot struct {
	Name            string             `yaml:"-"`
	Builder         *ZiplineeBuilder   `yaml:"builder,omitempty"`
	CloneRepository *bool              `yaml:"clone,omitempty" json:",omitempty"`
	Triggers        []*ZiplineeTrigger `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*ZiplineeStage   `yaml:"-" json:",omitempty"`
}

// UnmarshalYAML customizes unmarshalling an ZiplineeBot
func (bot *ZiplineeBot) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string             `yaml:"-"`
		Builder         *ZiplineeBuilder   `yaml:"builder"`
		CloneRepository *bool              `yaml:"clone"`
		Triggers        []*ZiplineeTrigger `yaml:"triggers"`
		Stages          yaml.MapSlice      `yaml:"stages"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	bot.Name = aux.Name
	bot.Builder = aux.Builder
	bot.CloneRepository = aux.CloneRepository
	bot.Triggers = aux.Triggers

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

		stage.Name = mi.Key.(string)
		bot.Stages = append(bot.Stages, stage)
	}

	return nil
}

// MarshalYAML customizes marshalling an ZiplineeBot
func (bot *ZiplineeBot) MarshalYAML() (out interface{}, err error) {

	var aux struct {
		Name            string             `yaml:"-"`
		Builder         *ZiplineeBuilder   `yaml:"builder,omitempty"`
		CloneRepository *bool              `yaml:"clone,omitempty"`
		Triggers        []*ZiplineeTrigger `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice      `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = bot.Builder
	aux.CloneRepository = bot.CloneRepository
	aux.Triggers = bot.Triggers

	for _, stage := range bot.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}
