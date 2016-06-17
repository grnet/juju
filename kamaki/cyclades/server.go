package cyclades

import (
	"fmt"

	"github.com/juju/juju/kamaki"
	"github.com/juju/juju/kamaki/client"
)

// Filter field keys.
const (
	NamePrefix = "--name-prefix"
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
	personality := formatPersonalityInfo(serverOpts.Personality)
	args := []string{"server", "create", "--name", serverOpts.Name,
		"--flavor-id", serverOpts.FlavorID, "--image-id", serverOpts.ImageID,
		"--project-id", serverOpts.ProjectID, "--output-format", "json", "-c",
		compute.Client.GetKamakirc()}
	args = append(args, personality...)
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
// Lists servers according to the filtering given as parameter.
// Returns a slice of the details of the servers or any error encountered.
func (compute Client) ListServers(filter map[string][]string) (
	[]ServerDetails, error) {
	filterValues := formatFilter(filter)
	args := []string{"server", "list", "-l", "--output-format",
		"json", "-c", compute.Client.GetKamakirc()}
	args = append(args, filterValues...)
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

// This function filtering options and formats them to a slice of strings.
// This format is required to pass filtering options as argument kamaki exec
// commands.
func formatFilter(filter map[string][]string) []string {
	var filterValues []string
	for k, filt := range filter {
		for _, v := range filt {
			filterValues = append(filterValues, []string{k, v}...)
		}
	}
	return filterValues
}

// This function takes a slice of personality info and formats them to a
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
