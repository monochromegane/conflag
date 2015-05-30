package conflag

import (
	"reflect"
	"testing"
)

func TestToArgs(t *testing.T) {
	c := conf{"flag": "value"}
	expect := []string{"-flag", "value"}
	actual := c.toArgs()

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
		actual := a.conf.toArgs(a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
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
		actual := a.conf.toArgs(a.positions...)
		if !reflect.DeepEqual(a.expect, actual) {
			t.Errorf("args should be %v, but %v", a.expect, actual)
		}
	}
}

type assert struct {
	conf      conf
	positions []string
	expect    []string
}
