package command

import "os/exec"

func HandleCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}
