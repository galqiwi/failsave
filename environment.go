package main

import (
	"os"
	"os/exec"
)

type Environment struct {
	Cmd     string            `json:"cmd"`
	CmdArgs []string          `json:"cmd_args"`
	EnvMap  map[string]string `json:"env_map"`
}

func (e *Environment) run() error {
	os.Clearenv()

	for envKey, envValue := range e.EnvMap {
		err := os.Setenv(envKey, envValue)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(e.Cmd, e.CmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
