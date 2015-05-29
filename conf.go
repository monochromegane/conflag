package conflag

import "fmt"

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
