package utils

import (
	"encoding/json"

	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"
)

// Dump JSON or YAML.
func Dump(config *viper.Viper, dumpType string) (string, error) {
	c := config.AllSettings()
	var bs []byte
	var err error
	switch dumpType {
	case "json":
		bs, err = json.Marshal(c)
	case "yaml":
		bs, err = yaml.Marshal(c)
	}
	if err != nil {
		return string(bs), err
	}
	return string(bs), nil
}
