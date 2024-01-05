package jump

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type JumpTargets struct {
	Targets map[string]string `json:"targets"`
}

func getHomDir() (homeDir string, err error) {
	homeDir, err = os.UserHomeDir()
	return
}

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

const jsonFileName = "targets.json"

func jsnFilePath() (path string) {

	homeDir, err := getHomDir()
	FatalOutError(err)
	configDir := ".config"
	jumpDir := "jump"
	configPath := filepath.Join(homeDir, configDir)
	jumpDirPath := filepath.Join(configPath, jumpDir)

	configNotExists := isPathExists(configPath)

	if configNotExists {
		err = os.Mkdir(configPath, os.FileMode(0755))
		FatalOutError(err)
	}

	jumpNotExists := isPathExists(jumpDirPath)

	if jumpNotExists {
		err = os.Mkdir(jumpDirPath, os.FileMode(0755))
		FatalOutError(err)
	}

	path = filepath.Join(jumpDirPath, jsonFileName)

	jsonNotExists := isPathExists(path)

	if jsonNotExists {
		initial := JumpTargets{
			Targets: make(map[string]string),
		}
		jsonContent, err := json.MarshalIndent(initial, "", "  ")
		FatalOutError(err)
		file, err := os.Create(path)
		FatalOutError(err)
		defer file.Close()
		_, err = file.Write(jsonContent)
		FatalOutError(err)
	}

	return
}

func jump(path string) {

	shellPath := os.Getenv("SHELL")
	_, shell := filepath.Split(shellPath)

	cmd := exec.Command(shell, "-c", fmt.Sprintf("cd %s && exec %s", path, shell))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func readTargets(path string) (jumpTargets JumpTargets) {
	bytes, err := os.ReadFile(path)
	FatalOutError(err)

	err = json.Unmarshal(bytes, &jumpTargets)
	FatalOutError(err)
	return
}

func getFullPath(path string) (fullPath string) {
	curDir, err := os.Getwd()
	FatalOutError(err)
	fullPath = filepath.Join(curDir, path)
	_, err = os.Stat(fullPath)
	FatalOutError(err)

	return
}

func ListTargets() {
	jsonPath := jsnFilePath()
	jumpTargets := readTargets(jsonPath)
	for key := range jumpTargets.Targets {
		fmt.Println(key)
	}
}

func DeleteTarget(target string) {

	jsonPath := jsnFilePath()
	jumpTargets := readTargets(jsonPath)
	delete(jumpTargets.Targets, target)
	targets, err := json.MarshalIndent(jumpTargets, "", " ")
	FatalOutError(err)
	err = os.WriteFile(jsonPath, targets, 0644)
	FatalOutError(err)

	fmt.Printf("%s is removed\n", target)
}

func AddTarget(path string) {

	fullPath := getFullPath(path)
	targetDir := filepath.Base(fullPath)

	jsonPath := jsnFilePath()

	jumpTargets := readTargets(jsonPath)
	jumpTargets.Targets[targetDir] = fullPath

	targets, err := json.MarshalIndent(jumpTargets, "", " ")
	FatalOutError(err)

	err = os.WriteFile(jsonPath, targets, 0644)
	FatalOutError(err)

	fmt.Printf("%s locked for jump\n", targetDir)

}

func InitiateJump(target string) {

	jsonPath := jsnFilePath()
	jumpTargets := readTargets(jsonPath)
	path, found := jumpTargets.Targets[target]

	if !found {
		fmt.Println("Target is not available to jump, you can add target by -a")
		os.Exit(1)
	}

	jump(path)
}
