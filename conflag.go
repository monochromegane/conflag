// Package conflag provides configuration parser for command-line flag.
package conflag

import "os"

// ArgsFrom make arguments for command-line flag from configuration file.
func ArgsFrom(conf string, positions ...string) ([]string, error) {
	if _, err := os.Stat(conf); err != nil {
		return nil, err
	}
	return parse(conf, positions...)
}
