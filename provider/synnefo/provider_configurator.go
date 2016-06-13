package synnefo

import (
	"bytes"
	"fmt"
	"github.com/juju/utils/ssh"
	"io"
	"strings"
)

func configureInstance(host string, login string, stdin io.Reader,
	stdout io.Writer) error {
	cmd := ssh.Command(login+"@"+host, []string{"sudo",
		"/bin/bash -c " + configureScript}, nil)
	var stderr bytes.Buffer
	cmd.Stdin = stdin
	cmd.Stdout = stdout // for sudo prompt
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() != 0 {
			err = fmt.Errorf("%v (%v)", err, strings.TrimSpace(stderr.String()))
		}
		return err
	}
	return nil
}

const configureScript = `
chown -R ubuntu:ubuntu /var/lib/juju
sudo add-apt-repository ppa:grnet/synnefo -y -s
sudo apt-get update -y
sudo apt-get install kamaki -y
`
