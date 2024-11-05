package manifest

import "time"

// ZiplineePipelineEvent fires for pipeline changes
type ZiplineePipelineEvent struct {
	BuildVersion string `yaml:"buildVersion,omitempty" json:"buildVersion,omitempty"`
	RepoSource   string `yaml:"repoSource,omitempty" json:"repoSource,omitempty"`
	RepoOwner    string `yaml:"repoOwner,omitempty" json:"repoOwner,omitempty"`
	RepoName     string `yaml:"repoName,omitempty" json:"repoName,omitempty"`
	Branch       string `yaml:"repoBranch,omitempty" json:"repoBranch,omitempty"`
	Status       string `yaml:"status,omitempty" json:"status,omitempty"`
	Event        string `yaml:"event,omitempty" json:"event,omitempty"`
}

// ZiplineeReleaseEvent fires for pipeline releases
type ZiplineeReleaseEvent struct {
	ReleaseVersion string `yaml:"releaseVersion,omitempty" json:"releaseVersion,omitempty"`
	RepoSource     string `yaml:"repoSource,omitempty" json:"repoSource,omitempty"`
	RepoOwner      string `yaml:"repoOwner,omitempty" json:"repoOwner,omitempty"`
	RepoName       string `yaml:"repoName,omitempty" json:"repoName,omitempty"`
	Target         string `yaml:"target,omitempty" json:"target,omitempty"`
	Status         string `yaml:"status,omitempty" json:"status,omitempty"`
	Event          string `yaml:"event,omitempty" json:"event,omitempty"`
}

// ZiplineeGitEvent fires for git repository changes
type ZiplineeGitEvent struct {
	Event      string `yaml:"event,omitempty" json:"event,omitempty"`
	Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
	Branch     string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// ZiplineeDockerEvent fires for docker image changes
type ZiplineeDockerEvent struct {
	Event string `yaml:"event,omitempty" json:"event,omitempty"`
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
	Tag   string `yaml:"tag,omitempty" json:"tag,omitempty"`
}

// ZiplineeCronEvent fires at intervals specified by the cron expression
type ZiplineeCronEvent struct {
	Time time.Time `yaml:"time,omitempty" json:"time,omitempty"`
}

// ZiplineeManualEvent fires when a user manually triggers a build or release
type ZiplineeManualEvent struct {
	UserID string `yaml:"userID,omitempty" json:"userID,omitempty"`
}

// ZiplineePubSubEvent fires when a subscribed pubsub topic receives an event
type ZiplineePubSubEvent struct {
	Project string        `yaml:"project,omitempty" json:"project,omitempty"`
	Topic   string        `yaml:"topic,omitempty" json:"topic,omitempty"`
	Message PubsubMessage `yaml:"message,omitempty" json:"message,omitempty"`
}

// ZiplineeGithubEvent fires for github events
type ZiplineeGithubEvent struct {
	Event      string `yaml:"event,omitempty" json:"event,omitempty"`
	Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
	Delivery   string `yaml:"delivery,omitempty" json:"delivery,omitempty"`
	Payload    string `yaml:"payload,omitempty" json:"payload,omitempty"`
}

// ZiplineeBitbucketEvent fires for bitbucket events
type ZiplineeBitbucketEvent struct {
	Event         string `yaml:"event,omitempty" json:"event,omitempty"`
	Repository    string `yaml:"repository,omitempty" json:"repository,omitempty"`
	HookUUID      string `yaml:"hookUUID,omitempty" json:"hookUUID,omitempty"`
	RequestUUID   string `yaml:"requestUUID,omitempty" json:"requestUUID,omitempty"`
	AttemptNumber string `yaml:"attemptNumber,omitempty" json:"attemptNumber,omitempty"`
	Payload       string `yaml:"payload,omitempty" json:"payload,omitempty"`
}

// ZiplineeEvent is a container for any trigger event
type ZiplineeEvent struct {
	Name      string                  `yaml:"name,omitempty" json:"name,omitempty"`
	Fired     bool                    `yaml:"fired,omitempty" json:"fired,omitempty"`
	Pipeline  *ZiplineePipelineEvent  `yaml:"pipeline,omitempty" json:"pipeline,omitempty"`
	Release   *ZiplineeReleaseEvent   `yaml:"release,omitempty" json:"release,omitempty"`
	Git       *ZiplineeGitEvent       `yaml:"git,omitempty" json:"git,omitempty"`
	Docker    *ZiplineeDockerEvent    `yaml:"docker,omitempty" json:"docker,omitempty"`
	Cron      *ZiplineeCronEvent      `yaml:"cron,omitempty" json:"cron,omitempty"`
	PubSub    *ZiplineePubSubEvent    `yaml:"pubsub,omitempty" json:"pubsub,omitempty"`
	Github    *ZiplineeGithubEvent    `yaml:"github,omitempty" json:"github,omitempty"`
	Bitbucket *ZiplineeBitbucketEvent `yaml:"bitbucket,omitempty" json:"bitbucket,omitempty"`
	Manual    *ZiplineeManualEvent    `yaml:"manual,omitempty" json:"manual,omitempty"`
}
