package conflag

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
)

func parse(file string, positions ...string) ([]string, error) {
	var conf conf

	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	switch filepath.Ext(file) {
	case ".toml":
		conf, err = parseAsToml(r)
	case ".json":
		conf, err = parseAsJson(r)
	case ".yaml", ".yml":
		conf, err = parseAsYaml(r)
	}
	if err != nil {
		return nil, err
	}
	return conf.toArgs(option{
		boolValue:  BoolValue,
		longHyphen: LongHyphen,
	}, positions...), nil
}

func parseAsToml(r io.Reader) (conf, error) {
	var conf conf
	_, err := toml.DecodeReader(r, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func parseAsJson(r io.Reader) (conf, error) {
	var conf conf
	if err := json.NewDecoder(r).Decode(&conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func parseAsYaml(r io.Reader) (conf, error) {
	var config conf

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
