package config

import (
	"io"
	"encoding/json"
	"io/ioutil"
	"gopkg.in/validator.v2"
	"github.com/pkg/errors"
)

type RotatorConfig struct {
	ActiveKeyLabel      string `json:"activeKeyLabel" validate:"nonzero"`
	ActiveKeyPassphrase string `json:"activeKeyPassphrase" validate:"nonzero"`
	DatabaseHostname    string `json:"databaseHostname" validate:"nonzero"`
	DatabasePort        string `json:"databasePort" validate:"nonzero"`
	DatabaseScheme      string `json:"databaseScheme" validate:"nonzero"`
	DatabaseName        string `json:"databaseName" validate:"nonzero"`
	DatabaseUsername    string `json:"databaseUsername" validate:"nonzero"`
	DatabasePassword    string `json:"databasePassword" validate:"nonzero"`
}

func New(rotatorConfigReader io.Reader) (*RotatorConfig, error) {

	rotatorConfigContent, err := ioutil.ReadAll(rotatorConfigReader)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read config")
	}

	rotatorConfig := &RotatorConfig{}
	err = json.Unmarshal(rotatorConfigContent, rotatorConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Malformed JSON provided.")
	}

	err = validator.Validate(rotatorConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid config.")
	}

	return rotatorConfig, nil
}
