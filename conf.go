package conflag

import "fmt"

type conf map[string]interface{}

type option struct {
	longHyphen bool
	boolValue  bool
}

func (c conf) toArgs(opt option, positions ...string) []string {

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
		for k, v := range n2 {
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
		hyphen := "-"
		if opt.longHyphen && len(k) > 1 {
			hyphen = "--"
		}
		name := hyphen + k

		switch v.(type) {
		case map[string]interface{}: // nested configuration
			continue
		case []interface{}:
			for _, a := range v.([]interface{}) {
				args = append(args, name, fmt.Sprintf("%v", a))
			}
		case bool:
			if opt.boolValue {
				args = append(args, fmt.Sprintf("%s=%t", name, v.(bool)))
			} else {
				if v.(bool) {
					args = append(args, name)
				}
			}
		default:
			args = append(args, name, fmt.Sprintf("%v", v))
		}
	}

	return args
}
