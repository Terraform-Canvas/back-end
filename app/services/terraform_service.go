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

var subnetList []models.Resource

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

func MergeEnvTf(tfFilePath string, userFolderPath string, resources []models.Resource) error {
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

	for _, item := range resources {
		folderPath := item.Type

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

			userVarPath := filepath.Join(userFolderPath, "variables.tf")
			if err := appendFile(userVarPath, varContent); err != nil {
				return err
			}

			if folderPath == "eks" {
				types_path := filepath.Join(tfFilePath, folderPath)
				var node_groups []byte

				eks_managed := item.Data["managed"].(map[string]interface{})
				if len(eks_managed) != 0 {
					contents, err := ioutil.ReadFile(filepath.Join(types_path, "eks_managed.txt"))
					node_groups = append(node_groups, contents...)
					vars, err := ioutil.ReadFile(filepath.Join(types_path, "eks_managed_var.txt"))
					err = appendFile(userVarPath, vars)
					if err != nil {
						return err
					}
				}
				fargate := item.Data["fargate"].(map[string]interface{})
				if len(fargate) != 0 {
					contents, err := ioutil.ReadFile(filepath.Join(types_path, "fargate.txt"))
					node_groups = append(node_groups, contents...)
					vars, err := ioutil.ReadFile(filepath.Join(types_path, "fargate_var.txt"))
					err = appendFile(userVarPath, vars)
					if err != nil {
						return err
					}
				}

				re := regexp.MustCompile(`//node_groups`)

				mainContent = []byte(re.ReplaceAllString(string(mainContent), string(node_groups)))
			}

			userMainPath := filepath.Join(userFolderPath, "main.tf")
			if err := appendFile(userMainPath, mainContent); err != nil {
				return err
			}

			re := regexp.MustCompile(`(?ms)module\s+"` + folderPath + `[^"]*"\s*{(?:[^{}]*{[^{}]*})*[^{}]*}`)
			if matches := re.FindAllString(string(sgFileContent), -1); len(matches) > 0 {
				sgContent := strings.Join(matches, "\n\n")
				if err := appendFile(userMainPath, []byte(sgContent)); err != nil {
					return err
				}
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

func CreateTfvars(userFolderPath string, resources []models.Resource) error {
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

	for _, item := range resources {
		itemType := item.Type
		itemData := item.Data

		// Process variables
		for key, value := range itemData {
			if name := fmt.Sprintf("%s_%s", itemType, key); variables[name] != nil {
				variables[name] = processValue(name, value)
			}
			if key == "fargate" || key == "managed" {
				for types, nval := range value.(map[string]interface{}) {
					if name := fmt.Sprintf("%s_%s_%s", itemType, key, types); variables[name] != nil {
						variables[name] = processValue(name, nval)
					}
				}
			}
		}

		// Process vpc subnet cidr
		if name := fmt.Sprintf("%s_publicsubnet", itemType); variables[name] != nil {
			subnetCnt := int(itemData["privatesubnet"].(float64) + itemData["publicsubnet"].(float64))
			req := models.SubnetRequest{
				VpcCidr:   itemData["cidr"].(string),
				SubnetCnt: subnetCnt,
			}
			res := calcSubnet(&req)

			variables[name] = res[:int(itemData["publicsubnet"].(float64))]
			variables["vpc_privatesubnet"] = res[int(itemData["publicsubnet"].(float64)):]
		}

		if name := fmt.Sprintf("%s_%s", itemType, "subnet_count"); variables[name] != nil {
			kind, start, end := subnetDepend(item) 
			vpcSubnet := variables[fmt.Sprintf("vpc_%s", kind)].([]string)
			if vpcSubnetLen := len(vpcSubnet); end > vpcSubnetLen {
				start -= len(subnetList) - vpcSubnetLen
				end -= len(subnetList) - vpcSubnetLen
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
		case int, int64, float64, bool:
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

func processValue(name string, value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		return toStringSlice(v)
	default:
		return v
	}
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

func subnetDepend(item models.Resource) (string, int, int) {
	tp, start, end := "", -1, -1
	for idx, sub := range subnetList {
		if sub.Parent == item.Id {
			if start == -1 {
				tp = sub.Type
				start, end = idx, idx+1
			} else {
				end += 1
			}
		} else if sub.Id == item.Parent {
			if start == -1 {
				tp = sub.Type
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
		cmd := exec.Command("bash", "-c", command)
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

	cmd := exec.Command("bash", "-c", "terraform destroy -auto-approve")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
