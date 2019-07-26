package urlshort

import (
	"reflect"
	"testing"
)

func TestParseYaml(t *testing.T) {
	input := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	expected := []pathURL{
		{Path: "/urlshort", URL: "https://github.com/gophercises/urlshort"},
		{Path: "/urlshort-final", URL: "https://github.com/gophercises/urlshort/tree/solution"},
	}

	um := urlMapping{}
	err := um.fromYaml([]byte(input))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(um.pu, expected) {
		t.Errorf("%v %v", um.pu, expected)
	}
}

func TestToMap(t *testing.T) {
	input := []pathURL{
		{Path: "/urlshort", URL: "https://github.com/gophercises/urlshort"},
		{Path: "/urlshort-final", URL: "https://github.com/gophercises/urlshort/tree/solution"},
	}

	expected := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	um := urlMapping{}
	um.pu = input

	result := um.toMap()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%v+, %v+", result, expected)
	}
}
