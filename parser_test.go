package conflag

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type ParseTest struct {
	configString string
	expected     map[string]interface{}
}

type ParseErrorTest struct {
	configString string
	expected     *regexp.Regexp
}

// helper
func (c conf) asMap() map[string]interface{} {
	return c
}

// JSON parser tests

func TestParseJson(t *testing.T) {
	asserts := []ParseTest{
		ParseTest{
			`{"flag":"value"}`,
			map[string]interface{}{"flag": "value"},
		},
		ParseTest{
			`{"flag1":"value1", "flag2":"value2"}`,
			map[string]interface{}{"flag1": "value1", "flag2": "value2"},
		},
	}

	for _, a := range asserts {
		reader := strings.NewReader(a.configString)

		actual, err := parseAsJson(reader)
		if err != nil {
			t.Errorf("parse error: %v", a.configString)
		}

		if !reflect.DeepEqual(a.expected, actual.asMap()) {
			t.Errorf("not match: %#v %#v", a.expected, actual.asMap())
		}
	}
}

func TestParseJson_ParseError(t *testing.T) {
	asserts := []ParseErrorTest{
		ParseErrorTest{`{'flag':'value'}`, regexp.MustCompile("invalid character '\\'' looking for beginning of object key string")},
	}

	for _, a := range asserts {
		reader := strings.NewReader(a.configString)
		_, err := parseAsJson(reader)
		if a.expected.MatchString(err.Error()) {
			t.Errorf("parse error: %v", err)
		}
	}
}

// YAML parser tests

func TestParseYaml(t *testing.T) {
	asserts := []ParseTest{
		ParseTest{
			`flag: value`,
			map[string]interface{}{"flag": "value"},
		},
		ParseTest{
			"flag1: value1\nflag2: value2\n",
			map[string]interface{}{"flag1": "value1", "flag2": "value2"},
		},
	}

	for _, a := range asserts {
		reader := strings.NewReader(a.configString)
		actual, err := parseAsYaml(reader)

		if err != nil {
			t.Errorf("parse error: %v %v", err, a.configString)
		}

		if !reflect.DeepEqual(a.expected, actual.asMap()) {
			t.Errorf("not match: %#v %#v", a.expected, actual.asMap())
		}
	}
}

func TestParseYaml_ParseError(t *testing.T) {
	asserts := []ParseErrorTest{
		ParseErrorTest{`flag - value`, regexp.MustCompile("yaml: unmarshal errors:")},
	}

	for _, a := range asserts {
		reader := strings.NewReader(a.configString)
		_, err := parseAsYaml(reader)
		if !a.expected.MatchString(err.Error()) {
			t.Errorf("expected error but")
		}
	}
}
