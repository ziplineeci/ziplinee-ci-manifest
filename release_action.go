package manifest

// ZiplineeReleaseAction represents an action on a release target that controls what happens by running the release stage
type ZiplineeReleaseAction struct {
	Name      string `yaml:"name" json:"name"`
	HideBadge bool   `yaml:"hideBadge,omitempty" json:"hideBadge,omitempty"`
}
