package services_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"main/app/models"
	"main/app/services"

	"github.com/stretchr/testify/assert"
)

func TestTerraformInitializeFolder(t *testing.T) {
	tmpDir := t.TempDir()

	err := services.InitializeFolder(tmpDir)
	assert.NoError(t, err)

	_, err = os.Stat(tmpDir)
	assert.NoError(t, err)
}

func TestTerraformMergeEnvTf(t *testing.T) {
	tmpDir := t.TempDir()
	resources := []models.Resource{
		{
			Id:   "resource_id",
			Type: "example_type",
			Data: map[string]interface{}{
				"example_variable": "new_value",
			},
		},
	}

	sgContent := []byte("module \"example_type_sg\" {\n  source={}\n  }")
	varContent := []byte("variable \"example_type_example_var\" {}\n")
	mainContent := []byte("module \"example_type\" {}\n")

	tfFilePath := filepath.Join(tmpDir, "platform", "terraform")
	userFolderPath := filepath.Join(tmpDir, "usertf")

	err := os.MkdirAll(filepath.Join(tfFilePath, "version"), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tfFilePath, "version", "versions.tf"), []byte("\x00"), 0o644)
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Join(tfFilePath, "example_type"), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tfFilePath, "example_type", "main.tf"), mainContent, 0o644)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tfFilePath, "example_type", "variables.tf"), varContent, 0o644)
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Join(tfFilePath, "sg"), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tfFilePath, "sg", "main.tf"), sgContent, 0o644)
	assert.NoError(t, err)

	err = os.MkdirAll(userFolderPath, os.ModePerm)
	assert.NoError(t, err)

	err = services.MergeEnvTf(tfFilePath, userFolderPath, resources)
	assert.NoError(t, err)

	mainContentMerged, err := ioutil.ReadFile(filepath.Join(userFolderPath, "main.tf"))
	assert.NoError(t, err)

	varContentMerged, err := ioutil.ReadFile(filepath.Join(userFolderPath, "variables.tf"))
	assert.NoError(t, err)

	assert.Contains(t, string(mainContentMerged), "module \"example_type\" {}")
	assert.Contains(t, string(mainContentMerged), "module \"example_type_sg\" {\n  source={}\n  }")
	assert.Contains(t, string(varContentMerged), "variable \"example_type_example_var\" {}")
}

func TestTerraformCreateTfvars(t *testing.T) {
	tmpDir := t.TempDir()

	variablesContent := []byte(`
	variable "example_type_example_variable" {
		type = string
		default = "default_value"
	}`)

	err := ioutil.WriteFile(filepath.Join(tmpDir, "variables.tf"), variablesContent, 0o644)
	assert.NoError(t, err)

	resources := []models.Resource{
		{
			Id:   "resource_id",
			Type: "example_type",
			Data: map[string]interface{}{
				"example_variable":   "new_value",
				"incorrect_variable": "new_value",
			},
		},
	}

	err = services.CreateTfvars(tmpDir, resources)
	assert.NoError(t, err)

	tfvarsContent, err := ioutil.ReadFile(filepath.Join(tmpDir, "terraform.tfvars"))
	assert.NoError(t, err)
	assert.Contains(t, string(tfvarsContent), "example_type_example_variable = \"new_value\"")
	assert.NotContains(t, string(tfvarsContent), "example_type_incorrect_variable = \"new_value\"")
}
