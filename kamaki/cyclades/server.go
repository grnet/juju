package cyclades

import (
	"fmt"

	"github.com/juju/juju/kamaki"
	"github.com/juju/juju/kamaki/client"
)

// This client executes kamaki commands for CRUD operations on Synnefo
// Virtual Servers.
type Client struct {
	client.Client
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
	// Custom Metadata of the server.
	Metadata map[string]string
}

// This struct defines all required information for creating a new Virtual
// Synnefo Server.
type ServerOpts struct {
	// Server Name.
	Name string
	// ID of project in which resources will be allocated.
	ProjectID string
	// Flavor ID.
	FlavorID string
	// Image ID of for the new Server.
	ImageID string
	// Slice of files to be injected to the new Server.
	Personality []PersonalityInfo
	// Custom Metadata for the server creation.
	Metadata map[string]string
	// Wait server to build.
	Wait bool
}

// This function creates a new Server to the Synnefo cloud specified by this
// client.
// It is created based on a specific image, flavor, name, assigned to a
// specific project and customized according the personality info given as
// parameter.
// Returns details of the created server or any error encountered.
func (compute Client) CreateServer(serverOpts ServerOpts) (
	*ServerDetails, error) {
	personality := FormatPersonalityInfo(serverOpts.Personality)
	metadata := FormatMetadata(serverOpts.Metadata)
	args := []string{"server", "create", "--name", serverOpts.Name,
		"--flavor-id", serverOpts.FlavorID, "--image-id", serverOpts.ImageID,
		"--project-id", serverOpts.ProjectID, "--output-format", "json", "-c",
		compute.Client.GetConfigFilePath()}
	args = append(args, personality...)
	args = append(args, metadata...)
	if serverOpts.Wait {
		args = append(args, "-w")
	}
	server := &ServerDetails{}
	output, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return nil, fmt.Errorf("Cannot create server")
	}
	kamaki.ToStruct(output, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// This functions lists servers of a specific Synnefo cloud according to this
// client.
// Returns a slice of the details of the servers or any error encountered.
func (compute Client) ListServers() (
	[]ServerDetails, error) {
	args := []string{"server", "list", "-l", "--output-format",
		"json", "-c", compute.Client.GetConfigFilePath()}
	output, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return nil, fmt.Errorf("Cannot list servers")
	}
	servers := make([]ServerDetails, 0)
	err = kamaki.ToStruct(output, &servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// This function takes metadata and formats them to a slice of strings.
// This format is required to pass custom metadta as an argument to the kamaki
// exec command.
func FormatMetadata(metadata map[string]string) []string {
	var formattedMetadata []string
	for k, v := range metadata {
		formattedMetadata = append(formattedMetadata,
			[]string{"-m", k + "=" + v}...)
	}
	return formattedMetadata
}

// This function takes a slice of personality info and formats them to a
// slice of strings.
// This format is required to pass personality info as argument to kamaki exec
// command.
func FormatPersonalityInfo(personalityInfo []PersonalityInfo) []string {
	var personality []string
	for _, data := range personalityInfo {
		personality = append(personality, []string{"-p",
			data.LocalPath + "," + data.RemotePath + "," + data.Owner +
				"," + data.Group + "," + data.Permission}...)
	}
	return personality
}
