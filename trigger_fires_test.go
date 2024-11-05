package manifest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZiplineePipelineTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfNegativeLookupBranchRegexDoesMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "!~ main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfNegativeLookupBranchRegexDoesNotMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "!~ main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfEventDoesNotMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "!= main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfStatusDoesNotMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfNameDoesNotMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-builder",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfBranchDoesNotMatch", func(t *testing.T) {

		event := ZiplineePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Branch: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestZiplineeReleaseTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := ZiplineeReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfEventDoesNotMatch", func(t *testing.T) {

		event := ZiplineeReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineeReleaseTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfStatusDoesNotMatch", func(t *testing.T) {

		event := ZiplineeReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfNameDoesNotMatch", func(t *testing.T) {

		event := ZiplineeReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-builder",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfBranchDoesNotMatch", func(t *testing.T) {

		event := ZiplineeReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "ziplineeci",
			RepoName:   "ziplinee-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-api",
			Target: "staging",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestZiplineeCronTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventTimeMatchesCronSchedule", func(t *testing.T) {

		event := ZiplineeCronEvent{
			Time: time.Date(2019, 4, 5, 11, 10, 0, 0, time.UTC),
		}

		trigger := ZiplineeCronTrigger{
			Schedule: "*/5 * * * *",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfEventTimeMatchesCronSchedule", func(t *testing.T) {

		event := ZiplineeCronEvent{
			Time: time.Date(2019, 4, 5, 11, 12, 1, 0, time.UTC),
		}

		trigger := ZiplineeCronTrigger{
			Schedule: "*/5 * * * *",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestZiplineeGitTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := ZiplineeGitEvent{
			Event:      "push",
			Repository: "bitbucket.org/xivart/icarus_to_email_service_trigger",
			Branch:     "main"}

		trigger := ZiplineeGitTrigger{
			Event:      "push",
			Repository: "bitbucket.org/xivart/icarus_to_email_service_trigger",
			Branch:     "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestZiplineeGithubTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventIsContainedInTriggerEvents", func(t *testing.T) {

		event := ZiplineeGithubEvent{
			Event: "create",
		}

		trigger := ZiplineeGithubTrigger{
			Events: []string{
				"commit_comment",
				"create",
				"delete",
				"deployment",
				"deployment_status",
				"fork",
				"gollum",
				"installation",
				"installation_repositories",
			},
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestZiplineeBitbucketTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventIsContainedInTriggerEvents", func(t *testing.T) {

		event := ZiplineeGithubEvent{
			Event: "pullrequest:comment_created",
		}

		trigger := ZiplineeGithubTrigger{
			Events: []string{
				"pullrequest:fulfilled",
				"pullrequest:rejected",
				"pullrequest:comment_created",
				"pullrequest:comment_updated",
				"pullrequest:comment_deleted",
			},
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestZiplineePubsubTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfTopicAndProjectMatch", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "my-project",
			Topic:   "my-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfProjectDoesNotMatch", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "another-project",
			Topic:   "my-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfTopicDoesNotMatch", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "my-project",
			Topic:   "another-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsTrueIfTopicAndProjectMatchAsRegex", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "my-project",
			Topic:   "my-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfProjectDoesNotMatchAsRegex", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "-project",
			Topic:   "my-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfTopicDoesNotMatchAsRegex", func(t *testing.T) {

		event := ZiplineePubSubEvent{
			Project: "my-project",
			Topic:   "-topic",
		}

		trigger := ZiplineePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}
