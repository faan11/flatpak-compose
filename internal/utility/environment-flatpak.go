package utility 

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"errors"
	"encoding/base64" 
	"io/ioutil"
	"github.com/faan11/flatpak-compose/internal/model"
)

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func ParseEnvironment(file_path string) (model.Environment, error) {
	var config model.Environment
	config.Core = make(map[string]string)
	config.Remotes = make(map[string]map[string]string)

	file, err := os.Open(file_path + "/config")
	if err != nil {
		return model.Environment{}, err
	}
	defer file.Close()


	var currentSectionValue string
	var currentSectionName string
	var remoteData map[string]string
	var newSectionValue string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			newSection := strings.TrimSpace(line[1 : len(line)-1])
			newSectionFields := strings.SplitN(newSection, " ", 2)
			newSectionName := newSectionFields[0]
			if (len(newSectionFields) == 2){
				newSectionValue = strings.Trim(newSectionFields[1],"\"")
			}

			if newSectionName == "remote" {
				if (currentSectionName == "remote") {
					config.Remotes[currentSectionValue] = remoteData
				}
				currentSectionName = newSectionName
				currentSectionValue = newSectionValue	
				remoteData = make(map[string]string)
				
				gpgKeyFilePath := file_path + "/"+ newSectionValue +".trustedkeys.gpg"	
				// Check gpg file existence
				if checkFileExists(gpgKeyFilePath) {
					// read the whole
					gpgData, err := ioutil.ReadFile(gpgKeyFilePath)
					if err != nil {
						return model.Environment{}, err
					}
					// Convert it to base64
					encodedString := base64.StdEncoding.EncodeToString(gpgData)
					// Assign it to GPGKey label
					remoteData["GPGKey"] = encodedString
				} else {
					remoteData["GPGKey"] = ""
				}
			} else {
				// it's core
				currentSectionName = newSectionName
			}
		} else {

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				return model.Environment{}, fmt.Errorf("malformed line: %s", line)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if currentSectionName == "core" {
				config.Core[key] = value
			} else if currentSectionName == "remote" {
				remoteData[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return model.Environment{}, err
	}

	if len(remoteData) > 0 {
		config.Remotes[currentSectionValue] = remoteData
	}

	return config, nil
}

func GetUserEnvironment() (model.Environment, error){
	repo_path := os.Getenv("HOME") + "/.local/share/flatpak/repo"
	config, err := ParseEnvironment(repo_path)
	if err != nil {
		return model.Environment{}, err
	}
	config.InstallationType = "user"
	return config, err
}

func GetSystemEnvironment() (model.Environment, error){
	repo_path := "/var/lib/flatpak/repo" 
	config, err := ParseEnvironment(repo_path)
	if err != nil {
		return model.Environment{}, err
	}
	config.InstallationType = "system"
	return config, err
}
/*
func main() {
	config, err := ParseFlatpakConfig(file_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("[core]")
	for k, v := range config.Core {
		fmt.Printf("%s=%s\n", k, v)
	}
	fmt.Println("")

	for key, remote := range config.Remotes {
		fmt.Printf("[remote \"%s\"]\n",key)
		for kr, vr := range remote {
			fmt.Printf("%s=%s\n", kr, vr)
		}
		fmt.Println("")
	}
}
*/
