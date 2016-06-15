package cyclades

import (
	"fmt"
	"github.com/juju/juju/kamaki"
)

// This client executes kamaki commands for CRUD operations on Synnefo
// Virtual Servers.
type Client struct {
	*kamaki.Client
}

// This struct represents personality info for creating Synnefo Servers.
// Typically, personality info defines a file to be injected to virtual
// servers file system.
type PersonalityInfo struct {
	// Local path of file to be injected.
	LocalPath string
	// Destination location inside server Image.
	RemotePath string
	// Owner of the injected file.
	Owner string
	// Group name of the injected file.
	Group string
	// Permission in octal of the injected file.
	Permission string
}

// This struct represents a Synnefo Virtual Server.
// It describes detailed information about it.
type ServerDetails struct {
	// Name of server.
	Name string
	// Server host.
	Host string `json:"SNF:fqdn"`
}

// This function creates a new Server to the Synnefo cloud specified by this
// client.
// It is created based on a specific image, flavor, name, assigned to a
// specific project and customized according the personality info given as
// parameter.
// Returns details of the created server or any error encountered.
func (compute Client) CreateServer(serverName string, projectID string,
	flavorID string, imageID string, personalityInfo []PersonalityInfo,
	wait bool) (
	*ServerDetails, error) {
	personality := formatPersonalityInfo(personalityInfo)
	args := []string{"server", "create", "--name", serverName,
		"--flavor-id", flavorID, "--image-id", imageID, "--project-id",
		projectID, "--output-format", "json", "-c",
		compute.Client.GetKamakirc()}
	args = append(args, personality...)
	if wait {
		args = append(args, "-w")
	}
	server := &ServerDetails{}
	output, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return nil, fmt.Errorf("Cannot create server: %s", string(output))
	}
	kamaki.ToStruct(output, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// This functions list servers of a specific Synnefo cloud according to this
// client.
// Lists servers whose name starts with the prefix given as parameter.
// Returns a slice of the details of the servers or any error encountered.
func (compute Client) ListServers(namePrefix string) (
	[]ServerDetails, error) {
	args := []string{"server", "list", "-l", "--name-prefix", namePrefix,
		"--output-format", "json", "-c", compute.Client.GetKamakirc()}
	output, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return nil, fmt.Errorf("Cannot list servers: %s", string(output))
	}
	servers := make([]ServerDetails, 0)
	err = kamaki.ToStruct(output, &servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// This function takes a slice of personality info and formatted them to a
// slice of strings.
// This format is required to pass personality info as argument to kamaki exec
// command.
func formatPersonalityInfo(personalityInfo []PersonalityInfo) []string {
	var personality []string
	for _, data := range personalityInfo {
		personality = append(personality, []string{"-p",
			data.LocalPath + "," + data.RemotePath + "," + data.Owner +
				"," + data.Group + "," + data.Permission}...)
	}
	return personality
}
