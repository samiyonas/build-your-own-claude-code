package llmtools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type BashArgs struct {
	Command string `json:"command"`
}

func Read(filePath string) (string, error) {
	filePath = strings.TrimPrefix(filePath, "/")

	clean := filepath.Clean(filePath)
	if strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	info, err := os.Stat(clean)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", fmt.Errorf("file path is a directory")
	}

	if strings.Contains(filePath, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func Write(filePath string, content string) error {
	filePath = strings.TrimPrefix(filePath, "/")

	clean := filepath.Clean(filePath)
	if strings.HasPrefix(clean, "..") {
		return fmt.Errorf("invalid file path")
	}

	if strings.Contains(filePath, "..") {
		return fmt.Errorf("invalid file path")
	}

	err := os.WriteFile(clean, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func Bash(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

func ReadTool(args string) (string, error) {
	var readArgs ReadArgs
	err := json.Unmarshal([]byte(args), &readArgs)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	content, err := Read(readArgs.FilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return content, nil
}

func WriteTool(args string) (string, error) {
	var writeArgs WriteArgs
	err := json.Unmarshal([]byte(args), &writeArgs)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	err = Write(writeArgs.FilePath, writeArgs.Content)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	return "File written successfully", nil
}

func BashTool(args string) (string, error) {
	var bashArgs BashArgs
	err := json.Unmarshal([]byte(args), &bashArgs)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	output, err := Bash(bashArgs.Command)
	if err != nil {
		return "", fmt.Errorf("error executing bash command: %v", err)
	}

	return output, nil
}
