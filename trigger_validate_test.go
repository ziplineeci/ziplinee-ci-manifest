package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZiplineeTriggerValidate(t *testing.T) {
	t.Run("ReturnsNoErrorIfOneTypeAndBuildActionIsSetForTriggerTypeBuild", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Event:  "finished",
				Status: "succeeded",
				Name:   "github.com/ziplineeci/ziplinee-ci-api",
				Branch: "master",
			},
			Git:    nil,
			Docker: nil,
			Cron:   nil,
			BuildAction: &ZiplineeTriggerBuildAction{
				Branch: "master",
			},
			ReleaseAction: nil,
		}

		// act
		err := trigger.Validate("build", "")

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorIfOneTypeAndReleaseActionIsSetForTriggerTypeRelease", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Event:  "finished",
				Status: "succeeded",
				Name:   "github.com/ziplineeci/ziplinee-ci-api",
				Branch: "master",
			},
			Git:         nil,
			Docker:      nil,
			Cron:        nil,
			BuildAction: nil,
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target: "development",
			},
		}

		// act
		err := trigger.Validate(TriggerTypeRelease, "development")

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorIfAllTypesAreEmpty", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: nil,
			Git:      nil,
			Docker:   nil,
			Cron:     nil,
			BuildAction: &ZiplineeTriggerBuildAction{
				Branch: "master",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target: "development",
			},
		}

		// act
		err := trigger.Validate("build", "")

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfAllActionsAreEmpty", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Event:  "finished",
				Status: "succeeded",
				Name:   "github.com/ziplineeci/ziplinee-ci-api",
				Branch: "master",
			},
			Git:           nil,
			Docker:        nil,
			Cron:          nil,
			BuildAction:   nil,
			ReleaseAction: nil,
		}

		// act
		err := trigger.Validate("build", "")

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfMoreThanOneTypeIsSet", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Event:  "finished",
				Status: "succeeded",
				Name:   "github.com/ziplineeci/ziplinee-ci-api",
				Branch: "master",
			},
			Git: &ZiplineeGitTrigger{
				Event:      "push",
				Repository: "github.com/ziplineeci/ziplinee-ci-builder",
				Branch:     "master",
			},
			Docker: nil,
			Cron:   nil,
			BuildAction: &ZiplineeTriggerBuildAction{
				Branch: "master",
			},
			ReleaseAction: nil,
		}

		// act
		err := trigger.Validate("build", "")

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfMoreThanOneActionIsSet", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Event:  "finished",
				Status: "succeeded",
				Name:   "github.com/ziplineeci/ziplinee-ci-api",
				Branch: "master",
			},
			Git:    nil,
			Docker: nil,
			Cron:   nil,
			BuildAction: &ZiplineeTriggerBuildAction{
				Branch: "master",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target: "development",
			},
		}

		// act
		err := trigger.Validate("build", "")

		assert.NotNil(t, err)
	})
}

func TestZiplineePipelineTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event:  "",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfStatusIsEmptyWhenEventIsFinished", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfStatusIsEmptyWhenEventIsStarted", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestZiplineeReleaseTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfStatusIsEmptyWhenEventIsFinished", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfStatusIsEmptyWhenEventIsStarted", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfTargetIsEmpty", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Target: "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/ziplineeci/ziplinee-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestZiplineeGitTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := ZiplineeGitTrigger{
			Event:      "",
			Repository: "github.com/ziplineeci/ziplinee-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfRepositoryIsEmpty", func(t *testing.T) {

		trigger := ZiplineeGitTrigger{
			Event:      "push",
			Repository: "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := ZiplineeGitTrigger{
			Event:      "push",
			Repository: "github.com/ziplineeci/ziplinee-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestZiplineeCronTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfScheduleIsEmpty", func(t *testing.T) {

		trigger := ZiplineeCronTrigger{
			Schedule: "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfScheduleIsInvalid", func(t *testing.T) {

		trigger := ZiplineeCronTrigger{
			Schedule: "0 * * * * *",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfScheduleIsValid", func(t *testing.T) {

		trigger := ZiplineeCronTrigger{
			Schedule: "*/5 * * * *",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}
