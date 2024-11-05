package manifest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalStage(t *testing.T) {
	t.Run("ReturnsUnmarshaledStage", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
shell: /bin/bash
workDir: /go/src/github.com/ziplineeci/ziplinee-ci-manifest
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish
when:
  server == 'ziplinee'`), &stage)

		assert.Nil(t, err)
		assert.Equal(t, "docker:17.03.0-ce", stage.ContainerImage)
		assert.Equal(t, "/bin/bash", stage.Shell)
		assert.Equal(t, "/go/src/github.com/ziplineeci/ziplinee-ci-manifest", stage.WorkingDirectory)
		assert.Equal(t, 2, len(stage.Commands))
		assert.Equal(t, "cp Dockerfile ./publish", stage.Commands[0])
		assert.Equal(t, "docker build -t ziplinee-ci-builder ./publish", stage.Commands[1])
		assert.Equal(t, "server == 'ziplinee'", stage.When)
	})

	t.Run("DefaultsShellToShIfNotPresent", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish
when:
  server == 'ziplinee'`), &stage)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "/bin/sh", stage.Shell)
	})

	t.Run("DefaultsShellToPowershellIfNotPresentForWindows", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish
when:
  server == 'ziplinee'`), &stage)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "windows",
		})

		assert.Nil(t, err)
		assert.Equal(t, "powershell", stage.Shell)
	})

	t.Run("DefaultsWhenToStatusEqualsSucceededIfNotPresent", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish`), &stage)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "status == 'succeeded'", stage.When)
	})

	t.Run("DefaultsWorkingDirectoryToZiplineeWorkIfNotPresent", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish`), &stage)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "/ziplinee-work", stage.WorkingDirectory)
	})

	t.Run("DefaultsWorkingDirectoryToZiplineeWorkIfNotPresentForWindows", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish`), &stage)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "windows",
		})

		assert.Nil(t, err)
		assert.Equal(t, "C:/ziplinee-work", stage.WorkingDirectory)
	})

	t.Run("ReturnsNonReservedSimplePropertyAsCustomProperty", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
unknownProperty1: value1
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish`), &stage)

		assert.Nil(t, err)
		assert.Equal(t, "value1", stage.CustomProperties["unknownProperty1"])
	})

	t.Run("ReturnsNonReservedArrayPropertyAsCustomProperty", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
unknownProperty3:
- supported1
- supported2
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish`), &stage)

		assert.Nil(t, err)
		assert.NotNil(t, stage.CustomProperties["unknownProperty3"])
		assert.Equal(t, "supported1", stage.CustomProperties["unknownProperty3"].([]interface{})[0].(string))
		assert.Equal(t, "supported2", stage.CustomProperties["unknownProperty3"].([]interface{})[1].(string))
	})
}

func TestJSONMarshalStage(t *testing.T) {
	t.Run("MarshalMapStringInterface", func(t *testing.T) {

		property := map[string]interface{}{
			"container": map[string]interface{}{
				"repository": "extension",
				"name":       "gke",
				"tag":        "alpha",
			},
		}

		// act
		bytes, err := json.Marshal(property)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"container\":{\"name\":\"gke\",\"repository\":\"extension\",\"tag\":\"alpha\"}}", string(bytes))
		}
	})

	t.Run("ReturnsMarshaledStageForNestedCustomProperties", func(t *testing.T) {

		var stage ZiplineeStage

		err := yaml.Unmarshal([]byte(`
image: extensions/gke:dev
container:
  repository: extensions`), &stage)

		assert.Nil(t, err)

		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
		})

		// act
		bytes, err := json.Marshal(stage)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"ContainerImage\":\"extensions/gke:dev\",\"Shell\":\"/bin/sh\",\"WorkingDirectory\":\"/ziplinee-work\",\"When\":\"status == 'succeeded'\",\"CustomProperties\":{\"container\":{\"repository\":\"extensions\"}}}", string(bytes))
		}
	})

	t.Run("ReturnsMarshaledStage", func(t *testing.T) {

		var stage ZiplineeStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
shell: /bin/bash
workDir: /go/src/github.com/ziplineeci/ziplinee-ci-manifest
commands:
- cp Dockerfile ./publish
- docker build -t ziplinee-ci-builder ./publish
when:
  server == 'ziplinee'`), &stage)

		// act
		bytes, err := json.Marshal(stage)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"ContainerImage\":\"docker:17.03.0-ce\",\"Shell\":\"/bin/bash\",\"WorkingDirectory\":\"/go/src/github.com/ziplineeci/ziplinee-ci-manifest\",\"Commands\":[\"cp Dockerfile ./publish\",\"docker build -t ziplinee-ci-builder ./publish\"],\"When\":\"server == 'ziplinee'\"}", string(bytes))
		}
	})

}

func TestValidateOnStage(t *testing.T) {
	t.Run("ReturnsErrorIfImageAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := ZiplineeStage{
			ContainerImage: "docker",
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfShellAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := ZiplineeStage{
			Shell: "/bin/sh",
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfWorkingDirectoryAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := ZiplineeStage{
			WorkingDirectory: "/ziplinee-work",
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfCommandsAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := ZiplineeStage{
			Commands: []string{"dotnet build"},
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfEnvvarsAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := ZiplineeStage{
			EnvVars: map[string]string{
				"ENVA": "value a",
			},
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfImageIsNotSetWithoutParallelStages", func(t *testing.T) {

		stage := ZiplineeStage{
			ContainerImage: "",
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfImageIsNotSetButHasService", func(t *testing.T) {

		stage := ZiplineeStage{
			ContainerImage: "",
			Services: []*ZiplineeService{
				&ZiplineeService{
					Name:           "cockroachdb",
					ContainerImage: "cockroachdb/cockroach:v19.2.0",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorWhenAllFieldsAreValid", func(t *testing.T) {

		stage := ZiplineeStage{
			ContainerImage: "docker",
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorWhenAllParallelStagesAreValid", func(t *testing.T) {

		stage := ZiplineeStage{
			ParallelStages: []*ZiplineeStage{
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&ZiplineeStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(ZiplineeBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})
}
