package conflag

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

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
	}
	if err != nil {
		return nil, err
	}
	return conf.toArgs(positions...), nil
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
