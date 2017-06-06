package mdqi

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v1"
)

type Conf struct {
	// Mdq contains configurations about mdq.
	Mdq ConfMdq `yaml:"mdq"`

	// Mdqi contains configurations about mdqi.
	Mdqi ConfMdqi `yaml:"mdqi"`
}

type ConfMdq struct {
	// Bin is the path to mdq command.
	Bin string `yaml:"bin"`

	// Config is the path to configuration file for mdq.
	Config string `yaml:"config"`
}

type ConfMdqi struct {
	// History is the path to history file.
	// Default value is "$HOME/.mdqi_history".
	History string `yaml:"history"`

	// DefaultTag is default value for mdq's --tag option.
	DefaultTag string `yaml:"default_tag"`

	// DefaultDisplay is default value for mdqi's /display command.
	DefaultDisplay string `yaml:"default_display"`
}

func ConfFromFile(path string) (Conf, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return Conf{}, errors.Wrap(err, "failed to read config file.")
	}

	var conf Conf
	if err := yaml.Unmarshal(body, &conf); err != nil {
		return Conf{}, errors.Wrap(err, "failed to unmarshal cofig file.")
	}

	return conf, nil
}
