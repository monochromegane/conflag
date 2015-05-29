package conflag

import "os"

func ArgsFrom(conf string, positions ...string) ([]string, error) {
	if _, err := os.Stat(conf); err != nil {
		return nil, err
	}
	return parse(conf, positions...)
}
