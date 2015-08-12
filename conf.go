package conflag

import "fmt"

type conf map[string]interface{}

func (c conf) toArgs(positions ...string) []string {

	nowConf := c
ConfLoop:
	for _, p := range positions {
		nextConf, ok := nowConf[p]
		if !ok {
			return []string{}
		}
		n, ok := nextConf.(map[string]interface{})
		if ok {
			nowConf = n
			continue
		}
		n2, ok := nextConf.(map[interface{}]interface{})
		if !ok {
			continue
		}
		// Convert map[interface{}] to map[string]
		tmp := map[string]interface{}{}
		for k, v := range(n2) {
			key, ok := k.(string)
			if !ok {
				continue ConfLoop
			}
			tmp[key] = v
		}
		nowConf = tmp
	}

	var args []string
	for k, v := range nowConf {
		switch v.(type) {
		case map[string]interface{}: // nested configuration
			continue
		case bool:
			args = append(args, fmt.Sprintf("-%s=%t", k, v.(bool)))
		default:
			args = append(args, "-"+k, fmt.Sprintf("%v", v))
		}
	}

	return args
}
