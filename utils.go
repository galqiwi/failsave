package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getCurrentExecutablePath() (string, error) {
	return os.Readlink("/proc/self/exe")
}

func encodeBase64(s string) string {
	output := strings.Builder{}

	e := base64.NewEncoder(base64.StdEncoding, &output)

	_, err := fmt.Fprint(e, s)
	if err != nil {
		panic(err)
	}

	err = e.Close()
	if err != nil {
		panic(err)
	}

	return output.String()
}

func encodeJson(v any) string {
	output, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(output)
}

func getCommandPath(cmdName string) (string, error) {
	cmd := exec.Command("which", cmdName)

	outputBuilder := strings.Builder{}
	cmd.Stdout = &outputBuilder

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := outputBuilder.String()

	if output == "" {
		return "", fmt.Errorf("command %s not found", cmdName)
	}

	if output[len(output)-1] != '\n' {
		return "", fmt.Errorf("which %s returned invalid result", cmdName)
	}

	output = output[:len(output)-1]

	return output, nil
}

func getShellPath() (string, error) {
	output, err := getCommandPath("bash")
	if err == nil {
		return output, nil
	}

	return getCommandPath("sh")
}
