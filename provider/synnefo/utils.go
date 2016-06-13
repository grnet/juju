package synnefo

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/juju/juju/juju/osenv"
)

const (
	nonceFile = "nonce.txt"
)

func writeNonceFile(nonce []byte) error {
	err := ioutil.WriteFile(path.Join(osenv.JujuXDGDataHomePath(), nonceFile),
		nonce, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}

func deleteNonceFile() error {
	err := os.Remove(path.Join(osenv.JujuXDGDataHomePath(), nonceFile))
	if err != nil {
		panic(err)
	}
	return nil
}
