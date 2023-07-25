package services

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"main/app/models"

	"github.com/hashicorp/hcl"
)

func InitializeFolder(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func MergeEnvTf(userFolderPath string, data []map[string]interface{}) error {
	tfFilePath := filepath.Join("platform", "terraform")
	sgFilePath := filepath.Join(tfFilePath, "sg")
	sgFileContent, err := ioutil.ReadFile(sgFilePath)
	if err != nil {
		return err
	}

	for _, item := range data {
		folderPath := item["type"].(string)

		mainFilePath := filepath.Join(tfFilePath, folderPath, "main.tf")
		mainContent, err := ioutil.ReadFile(mainFilePath)
		if err != nil {
			return err
		}

		varFilePath := filepath.Join(tfFilePath, folderPath, "variables.tf")
		varContent, err := ioutil.ReadFile(varFilePath)
		if err != nil {
			return err
		}

		userMainPath := filepath.Join(userFolderPath, "main.tf")
		if err := AppendFile(userMainPath, mainContent); err != nil {
			return err
		}

		re := regexp.MustCompile(`(?ms)module\s+"` + folderPath + `[^"]*"\s*{(?:[^{}]*{[^{}]*})*[^{}]*}`)
		matches := re.FindAllString(string(sgFileContent), -1)

		if len(matches) > 0 {
			sgContent := strings.Join(matches, "\n\n")
			if err := AppendFile(userMainPath, []byte(sgContent)); err != nil {
				return err
			}
		}
		userVarPath := filepath.Join(userFolderPath, "variables.tf")
		if err := AppendFile(userVarPath, varContent); err != nil {
			return err
		}
	}

	return nil
}

func AppendFile(filePath string, content []byte) error {
	exist, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ioutil.WriteFile(filePath, content, 0o644)
	}

	newContent := append(exist, []byte("\n\n")...)
	newContent = append(newContent, content...)

	err = ioutil.WriteFile(filePath, newContent, 0o644)
	if err != nil {
		return nil
	}

	return nil
}

func CreateTfvars(userFolderPath string, data []map[string]interface{}) error {
	var v map[string]interface{}
	variablesFile := filepath.Join(userFolderPath, "variables.tf")
	existingVariables, err := ioutil.ReadFile(variablesFile)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`\{[^{}]*\}`)
	fileContent := re.ReplaceAllString(string(existingVariables), "{}")

	err = hcl.Unmarshal([]byte(fileContent), &v)
	if err != nil {
		return err
	}

	variableSlice := v["variable"].([]map[string]interface{})
	variables := make(map[string]interface{})
	for _, variable := range variableSlice {
		for key, value := range variable {
			variables[key] = value
		}
	}

	for _, item := range data {
		itemType := item["type"].(string)
		itemData := item["data"].(map[string]interface{})

		// prefix가 'itemType_'인 variables 항목을 찾아야 함.
		for name := range variables {
			if strings.HasPrefix(name, fmt.Sprintf("%s_", itemType)) {
				for key, value := range itemData {
					if name == "vpc_privatesubnet" || name == "vpc_publicsubnet" {
						subnetCnt := int(itemData["privatesubnet"].(float64) + itemData["publicsubnet"].(float64))

						req := models.SubnetRequest{
							VpcCidr:   itemData["cidr"].(string),
							SubnetCnt: subnetCnt,
						}

						res := CalcSubnet(&req)

						if name == "vpc_privatesubnet" {
							variables[name] = res[:int(itemData["privatesubnet"].(float64))]
						} else {
							variables[name] = res[int(itemData["privatesubnet"].(float64)):]
						}
						continue
					}
					if name == fmt.Sprintf("%s_%s", itemType, key) {
						variables[name] = value
					}
				}
			}
		}
	}

	var tfvars strings.Builder

	for key, value := range variables {
		tfvars.WriteString(fmt.Sprintf("%s = ", key))
		switch v := value.(type) {
		case string:
			tfvars.WriteString(fmt.Sprintf(`"%s"`, v))
		case int, int64, float64:
			tfvars.WriteString(fmt.Sprintf("%v", v))
		case []string:
			tfvars.WriteString("[")
			for i, item := range v {
				tfvars.WriteString(fmt.Sprintf(`"%s"`, item))
				if i != len(v)-1 {
					tfvars.WriteString(", ")
				}
			}
			tfvars.WriteString("]")
		default:
			continue
		}
		tfvars.WriteString("\n")
	}

	writePath := filepath.Join(userFolderPath, "terraform.tfvars")
	err = ioutil.WriteFile(writePath, []byte(tfvars.String()), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func CalcSubnet(req *models.SubnetRequest) []string {
	vpcCidr := req.VpcCidr
	subnetCnt := req.SubnetCnt

	parts := strings.Split(vpcCidr, "/")
	ip := parts[0]
	prefix := parts[1]

	ipParts := strings.Split(ip, ".")
	var ipInt uint32
	for _, part := range ipParts {
		num, err := strconv.Atoi(part)
		if num < 0 || num > 255 || err != nil {
			return nil
		}
		ipInt = (ipInt << 8) | uint32(num)
	}

	prefixLen, err := strconv.Atoi(prefix)
	if prefixLen < 0 || prefixLen > 32 || err != nil {
		return nil
	}

	subnetBits := uint32(math.Ceil(math.Log2(float64(subnetCnt))))
	subnetCidrs := make([]string, subnetCnt)
	for i := 0; i < subnetCnt; i++ {
		subnetIp := ipInt | uint32(i)<<(32-prefixLen-int(subnetBits))
		subnetCidr := fmt.Sprintf("%d.%d.%d.%d/%d", (subnetIp>>24)&255, (subnetIp>>16)&255, (subnetIp>>8)&255, subnetIp&255, prefixLen+int(subnetBits))
		subnetCidrs[i] = subnetCidr
	}

	return subnetCidrs
}

func ApplyTerraform(userFolderPath string) error {
	err := os.Chdir(userFolderPath)
	if err != nil {
		return err
	}

	commands := []string{
		"terraform fmt",
		"terraform init",
		"terraform validate",
		"terraform plan",
		"terraform apply -auto-approve",
	}

	for _, command := range commands {
		cmd := exec.Command(command)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func DestroyTerraform(userFolderPath string) error {
	err := os.Chdir(userFolderPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("terraform destroy -auto-approve")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
