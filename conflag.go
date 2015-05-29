package conflag

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func ArgsFrom(conf string, positions ...string) ([]string, error) {
	if _, err := os.Stat(conf); err == nil {
		return nil, err
	}
	return parse(conf, positions...)
}

func parse(file string, positions ...string) ([]string, error) {
	var conf conf
	var err error
	switch filepath.Ext(file) {
	case ".toml":
		conf, err = parseAsToml(file)
	}
	if err != nil {
		return nil, err
	}
	return conf.toArgs(positions...), nil
}

func parseAsToml(file string) (conf, error) {
	var conf conf
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

type conf map[string]interface{}

func (c conf) toArgs(positions ...string) []string {

	nowConf := c
	for _, p := range positions {
		nextConf, ok := nowConf[p]
		if !ok {
			break
		}
		n, ok := nextConf.(map[string]interface{})
		if !ok {
			break
		}
		nowConf = n
	}

	var args []string
	for k, v := range nowConf {
		args = append(args, "-"+k, fmt.Sprintf("%v", v))
	}

	return args
}
