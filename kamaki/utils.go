package kamaki

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"os/user"
	"path"
)

const (
	command = "kamaki"
)

// Starts the kamaki command with the given arguments and waits to be completed.
// It returns the standard output and any encounted error.
func RunCmdOutput(args []string) ([]byte, error) {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// This function gets the absolute path of kamakirc file located in the home
// directory of the current user.
// It returns the absolute path of kamakirc file or any error encountered.
func FormatPath(kamakirc string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Cannot extract home dir")
	}
	return path.Join(usr.HomeDir, kamakirc), nil
}

// This function converts raw JSON data to the struct given as parameter.
// It returns any error encountered.
func ToStruct(data []byte, structType interface{}) error {
	return json.Unmarshal(data, structType)
}
