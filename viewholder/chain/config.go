package chain

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Consensus []Node `yaml:"consensus"`
	// PrivateKey string `yaml:"private_key"`
	// SignCert   string `yaml:"sign_cert"`
	Singers           []Node `yaml:"singers"`
	CLTMST            int    `yaml:"cltmsSignT"`
	CLTMSD            int    `yaml:"cltmsSignD"`
	UserChaincodeName string `yaml:"user_chaincode_name"`
}

type Node struct {
	TLSCACert     string `yaml:"tls_ca_cert"`
	TLSCACertByte []byte
	Domain        string `yaml:"domain"`
}

func (n *Node) loadConfig() error {
	TLSCACert, err := GetTLSCACerts(n.TLSCACert)
	if err != nil {
		return errors.Wrapf(err, "fail to load TLS CA Cert %s", n.TLSCACert)
	}
	n.TLSCACertByte = TLSCACert
	return nil
}

func GetTLSCACerts(file string) ([]byte, error) {
	if len(file) == 0 {
		return nil, nil
	}
	in, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, errors.Wrapf(err, "error loading %s", file)
	}
	return in, nil
}

func LoadConfig(file string) (Config, error) {
	config := Config{}

	readFile, err := ioutil.ReadFile(file)
	if err != nil {
		return config, errors.Wrapf(err, "err loading %s", file)
	}

	err = yaml.Unmarshal(readFile, &config)
	if err != nil {
		return config, errors.Wrapf(err, "error unmarshal %s", file)
	}

	for _, consensus := range config.Consensus {
		err = consensus.loadConfig()
		if err != nil {
			return config, err
		}
	}

	return config, nil

}
