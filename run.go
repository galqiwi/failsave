package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var recoveryScriptPath string

var RunCmd = &cobra.Command{
	Short: "run command and save environment",
	Use:   "run [flags] command [command args]",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args)
	},
}

func init() {
	RunCmd.PersistentFlags().StringVar(
		&recoveryScriptPath,
		"recovery-path",
		"./recovery.sh",
		"Recovery script path",
	)
}

func BuildRecoveryScript(env *Environment) (string, error) {
	output := strings.Builder{}

	shellPath, err := getShellPath()
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprintf(&output, "#!%s\n\n", shellPath)
	if err != nil {
		return "", err
	}

	selfPath, err := getCurrentExecutablePath()
	if err != nil {
		return "", err
	}

	// sanitization through base64
	_, err = fmt.Fprintf(
		&output,
		"\"$(echo %s | base64 -d)\" recovery %s\n",
		encodeBase64(selfPath),
		encodeBase64(encodeJson(env)),
	)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func getEnvMap() map[string]string {
	output := make(map[string]string)

	for _, entry := range os.Environ() {
		entrySplitted := strings.Split(entry, "=")
		if len(entrySplitted) < 1 {
			panic("invalid os.Environ call")
		}
		envName := entrySplitted[0]

		envValue := os.Getenv(envName)
		output[envName] = envValue
	}

	return output
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command specified")
	}
	cmd := args[0]
	args = args[1:]

	env := &Environment{
		Cmd:     cmd,
		CmdArgs: args,
		EnvMap:  getEnvMap(),
	}

	recoveryScript, err := BuildRecoveryScript(env)
	if err != nil {
		return err
	}

	err = os.WriteFile(recoveryScriptPath, []byte(recoveryScript), 0777)
	if err != nil {
		return err
	}

	return env.run()
}
