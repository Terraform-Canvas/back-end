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

var subnetList []map[string]interface{}

func InitializeFolder(folderPath string) error {
	subnetList = subnetList[:0]
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

	// version config
	versionPath := filepath.Join(tfFilePath, "version", "versions.tf")
	versionContent, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return err
	}
	err = createFile(userFolderPath, "versions.tf", versionContent)
	if err != nil {
		return err
	}

	// sg config
	sgFilePath := filepath.Join(tfFilePath, "sg", "main.tf")
	sgFileContent, err := ioutil.ReadFile(sgFilePath)
	if err != nil {
		return err
	}

	for _, item := range data {
		folderPath := item["type"].(string)

		if folderPath != "privatesubnet" && folderPath != "publicsubnet" {
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
			if err := appendFile(userMainPath, mainContent); err != nil {
				return err
			}

			re := regexp.MustCompile(`(?ms)resource\s+"[^"]*"\s"` + folderPath + `[^"]*"\s*{(?:[^{}]*{[^{}]*})*[^{}]*}`)
			if matches := re.FindAllString(string(sgFileContent), -1); len(matches) > 0 {
				sgContent := strings.Join(matches, "\n\n")
				if err := appendFile(userMainPath, []byte(sgContent)); err != nil {
					return err
				}
			}

			userVarPath := filepath.Join(userFolderPath, "variables.tf")
			if err := appendFile(userVarPath, varContent); err != nil {
				return err
			}
		} else {
			subnetList = append(subnetList, item)
		}
	}

	return nil
}

func appendFile(filePath string, content []byte) error {
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

		// Process variables
		for key, value := range itemData {
			if name := fmt.Sprintf("%s_%s", itemType, key); variables[name] != nil {
				if i, ok := value.([]interface{}); ok {
					variables[name] = toStringSlice(i)
				} else {
					variables[name] = value
				}
			}
		}

		// Process vpc subnet cidr
		if name := fmt.Sprintf("%s_privatesubnet", itemType); variables[name] != nil {
			subnetCnt := int(itemData["privatesubnet"].(float64) + itemData["publicsubnet"].(float64))
			req := models.SubnetRequest{
				VpcCidr:   itemData["cidr"].(string),
				SubnetCnt: subnetCnt,
			}
			res := calcSubnet(&req)

			variables[name] = res[:int(itemData["privatesubnet"].(float64))]
			variables["vpc_publicsubnet"] = res[int(itemData["privatesubnet"].(float64)):]
		}

		if name := fmt.Sprintf("%s_%s", itemType, "subnet_count"); variables[name] != nil {
			kind, start, end := subnetDepend(item)
			vpcSubnet := variables[fmt.Sprintf("vpc_%s", kind)].([]string)
			if vpcSubnetLen := len(vpcSubnet); end > vpcSubnetLen {
				start %= vpcSubnetLen
				end %= vpcSubnetLen
			}
			variables[name] = []int{start, end}
			variables[fmt.Sprintf("%s_%s", itemType, "subnet_type")] = kind
		}

		// Process user-data
		if itemData["user_data"] != nil {
			err = createFile(userFolderPath, "user-data.sh", []byte(itemData["user_data"].(string)))
			if err != nil {
				return err
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
		case []int:
			tfvars.WriteString("[")
			for i, item := range v {
				tfvars.WriteString(fmt.Sprintf("%v", item))
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

func toStringSlice(slice []interface{}) []string {
	result := make([]string, 0)
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func calcSubnet(req *models.SubnetRequest) []string {
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

func subnetDepend(item map[string]interface{}) (string, int, int) {
	tp, start, end := "", -1, -1
	for idx, sub := range subnetList {
		if sub["parent"] == item["id"].(string) {
			if start == -1 {
				tp = sub["type"].(string)
				start, end = idx, idx+1
			} else {
				end += 1
			}
		} else if sub["id"] == item["parent"].(string) {
			if start == -1 {
				tp = sub["type"].(string)
				start, end = idx, idx+1
			} else {
				end += 1
			}
		}
	}
	return tp, start, end
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

func createFile(userFolderPath string, fileName string, content []byte) error {
	filePath := filepath.Join(userFolderPath, fileName)

	err := ioutil.WriteFile(filePath, content, 0o644)
	if err != nil {
		return err
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
