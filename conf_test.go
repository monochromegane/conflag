package conflag

import (
	"bytes"
	"reflect"
	"testing"
)

func defaultOption() option {
	return option{
		boolValue:  true,
		longHyphen: false,
	}
}

func TestToArgs(t *testing.T) {
	c := conf{"flag": "value"}
	expect := []string{"-flag", "value"}
	actual := c.toArgs(defaultOption())

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("args should be %v, but %v", expect, actual)
	}
}

func TestToArgs_Bool(t *testing.T) {
	asserts := []assert{
		assert{conf{"flag": false}, nil, []string{"-flag=false"}},
		assert{conf{"flag": true}, nil, []string{"-flag=true"}},
	}

	for _, a := range asserts {
		actual := a.conf.toArgs(defaultOption(), a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

func TestToArgs_List(t *testing.T) {
	c := conf{"flag": []interface{}{"value1", "value2"}}
	expect := []string{"-flag", "value1", "-flag", "value2"}
	actual := c.toArgs(defaultOption())

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("args should be %v, but %v", expect, actual)
	}
}

func TestToArgs_Positions(t *testing.T) {
	asserts := []assert{
		// options in options section.
		assert{
			conf{"options": map[string]interface{}{"flag": "value"}},
			[]string{"options"},
			[]string{"-flag", "value"},
		},
		// options in general/options section.
		assert{
			conf{"general": map[string]interface{}{"options": map[string]interface{}{"flag": "value"}}},
			[]string{"general", "options"},
			[]string{"-flag", "value"},
		},
	}

	for _, a := range asserts {
		actual := a.conf.toArgs(defaultOption(), a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

func TestToArgs_LongHyphenOption(t *testing.T) {
	asserts := []assert{
		assert{conf{"f": "1"}, nil, []string{"-f", "1"}},
		assert{conf{"flag": "1"}, nil, []string{"--flag", "1"}},
	}

	for _, a := range asserts {
		actual := a.conf.toArgs(option{longHyphen: true}, a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

func TestToArgs_BoolValueOption(t *testing.T) {
	asserts := []assert{
		assert{conf{"flag": false}, nil, nil},
		assert{conf{"flag": true}, nil, []string{"-flag"}},
	}

	for _, a := range asserts {
		actual := a.conf.toArgs(option{boolValue: false}, a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

func TestToArgs_Positions2(t *testing.T) {
	asserts := []assert{
		assert{
			conf{
				"env1": map[interface{}]interface{}{"options": map[interface{}]interface{}{"flag": "value"}},
				"env2": map[interface{}]interface{}{"options": map[interface{}]interface{}{"flag": 123}},
			},
			[]string{"env2", "options"},
			[]string{"-flag", "123"},
		},
	}
	for _, a := range asserts {
		actual := a.conf.toArgs(defaultOption(), a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

func TestTurnAround_Toml(t *testing.T) {
	src := bytes.NewReader([]byte(`[env1]
flag1 = "value1"
flag2 = "value2"
[env2]
flag1 = "12345"
flag2 = "-1.41421356"`))
	conf, _ := parseAsToml(src)

	actual := conf.toArgs(defaultOption(), "env2")

	expect := []string{"-flag1", "12345", "-flag2", "-1.41421356"}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("args should be %v, but %v", expect, actual)
	}
}

func TestTurnAround_Yaml(t *testing.T) {
	src := bytes.NewReader([]byte(`env1:
  flag1: value1
  flag2: "value2"
# line comment
env2:
  flag1: 12345	# inline comment
  flag2: -1.41421356`))
	conf, _ := parseAsYaml(src)

	actual := conf.toArgs(defaultOption(), "env2")

	expect := []string{"-flag1", "12345", "-flag2", "-1.41421356"}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("args should be %v, but %v", expect, actual)
	}
}

type assert struct {
	conf      conf
	positions []string
	expect    []string
}
