package conflag

import (
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type ParseTest struct {
	configString string
	expected     interface{}
}

type ParseErrorTest struct {
	configString string
	expected     string // regexp
}

type parseFunc func(io.Reader) (conf, error)

// helper
func (c conf) asMap() map[string]interface{} {
	return c
}

//
// Common assert logic
//

func assertTestParse(t *testing.T, p parseFunc, testCase ParseTest) {
	reader := strings.NewReader(testCase.configString)
	actual, err := p(reader)

	if err != nil {
		t.Errorf("Unexpected parse error: %v, for input %#v", err, testCase.configString)
	}

	if !reflect.DeepEqual(testCase.expected, actual.asMap()) {
		t.Errorf("Parsed result should be %#v, but %#v", testCase.expected, actual.asMap())
	}
}

func assertTestParseError(t *testing.T, p parseFunc, testCase ParseErrorTest) {
	reader := strings.NewReader(testCase.configString)
	_, err := p(reader)

	re := regexp.MustCompile(testCase.expected)
	if !re.MatchString(err.Error()) {
		t.Errorf("Unexpected parse error: %v", err)
	}
}

//
// JSON parser tests
//

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
		assertTestParse(t, parseAsJson, a)
	}
}

func TestParseJson_ParseError(t *testing.T) {
	asserts := []ParseErrorTest{
		ParseErrorTest{`{'flag':'value'}`, `invalid character '\\'' looking for beginning of object key string`},
	}

	for _, a := range asserts {
		assertTestParseError(t, parseAsJson, a)
	}
}

//
// TOML parser tests
//

func TestParseToml(t *testing.T) {
	asserts := []ParseTest{
		ParseTest{
			`flag = "value"`,
			map[string]interface{}{"flag": "value"},
		},
		ParseTest{
			`flag1 = "value1"` + "\n" + `flag2 = "value2"`,
			map[string]interface{}{"flag1": "value1", "flag2": "value2"},
		},
	}

	for _, a := range asserts {
		assertTestParse(t, parseAsToml, a)
	}
}

func TestParseToml_ParseError(t *testing.T) {
	asserts := []ParseErrorTest{
		ParseErrorTest{`flag : value`, `Expected key separator '=', but got ':' instead`},
	}

	for _, a := range asserts {
		assertTestParseError(t, parseAsToml, a)
	}
}

//
// YAML parser tests
//

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
		ParseTest{
`env1:
  flag1: value1
  flag2: "value2"
# line comment
env2:
  flag1: 12345	# inline comment
  flag2: -1.41421356`,
			interface{}(map[string]interface{}{
				"env1":map[interface{}]interface{}{"flag1": "value1", "flag2": "value2"},
				"env2":map[interface{}]interface{}{"flag1": 12345, "flag2": -1.41421356},
			}),
		},

	}

	for _, a := range asserts {
		assertTestParse(t, parseAsYaml, a)
	}
}

func TestParseYaml_ParseError(t *testing.T) {
	asserts := []ParseErrorTest{
		ParseErrorTest{`flag - value`, "yaml: unmarshal errors:"},
	}

	for _, a := range asserts {
		assertTestParseError(t, parseAsYaml, a)
	}
}
