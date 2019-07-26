package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	ret := func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.RequestURI]; ok {
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return ret
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type pathURL struct {
	Path string
	URL  string
}

type urlMapping struct {
	pu []pathURL
}

func (um *urlMapping) fromYaml(yml []byte) error {
	err := yaml.Unmarshal(yml, &um.pu)
	if err != nil {
		return err
	}
	return nil
}

func (um *urlMapping) toMap() map[string]string {
	pathsToUrls := make(map[string]string)
	for _, urlMapping := range um.pu {
		pathsToUrls[urlMapping.Path] = urlMapping.URL
	}
	return pathsToUrls
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	um := urlMapping{}

	err := um.fromYaml(yml)
	if err != nil {
		return nil, err
	}

	// Create map from array
	pathsToUrls := um.toMap()
	ret := MapHandler(pathsToUrls, fallback)
	return ret, nil
}
