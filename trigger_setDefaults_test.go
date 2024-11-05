package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZiplineePipelineTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToFinishedIfEmpty", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "finished", trigger.Event)
	})

	t.Run("SetsStatusToSucceededIfEmpty", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Status: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Status)
	})

	t.Run("SetsBranchToMasterOrMainIfEmpty", func(t *testing.T) {

		trigger := ZiplineePipelineTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master|main", trigger.Branch)
	})
}

func TestZiplineeReleaseTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToFinishedIfEmpty", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "finished", trigger.Event)
	})
	t.Run("SetsStatusToSucceededIfEmpty", func(t *testing.T) {

		trigger := ZiplineeReleaseTrigger{
			Status: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Status)
	})
}

func TestZiplineeGitTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToPushIfEmpty", func(t *testing.T) {

		trigger := ZiplineeGitTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "push", trigger.Event)
	})

	t.Run("SetsBranchToMasterOrMainIfEmpty", func(t *testing.T) {

		trigger := ZiplineeGitTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master|main", trigger.Branch)
	})
}

func TestZiplineeTriggerBuildActionSetDefaults(t *testing.T) {
	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := ZiplineeTriggerBuildAction{
			Branch: "",
		}

		preferences := ZiplineeManifestPreferences{
			DefaultBranch: "main",
		}

		// act
		trigger.SetDefaults(preferences)

		assert.Equal(t, "main", trigger.Branch)
	})
}

func TestZiplineeTriggerReleaseActionSetDefaults(t *testing.T) {
	t.Run("SetsTargetToTargetParam", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target: "any",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "development", trigger.ReleaseAction.Target)
	})

	t.Run("SetsVersionToLatestIfEmptyAndTriggersOnOtherPipeline", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Name: "github.com/ziplineeci/ziplinee-ci-builder",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target:  "any",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "latest", trigger.ReleaseAction.Version)
	})

	t.Run("KeepsVersionIfNotEmpty", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target:  "any",
				Version: "current",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "current", trigger.ReleaseAction.Version)
	})

	t.Run("SetsVersionToSameIfPipelineTriggerIsTheSelfPipeline", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Pipeline: &ZiplineePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target:  "development",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "same", trigger.ReleaseAction.Version)
	})

	t.Run("SetsVersionToSameIfReleaseTriggerIsTheSelfPipeline", func(t *testing.T) {

		trigger := ZiplineeTrigger{
			Release: &ZiplineeReleaseTrigger{
				Name: "self",
			},
			ReleaseAction: &ZiplineeTriggerReleaseAction{
				Target:  "development",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "same", trigger.ReleaseAction.Version)
	})

}
