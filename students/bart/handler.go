package main

import (
	"fmt"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc that will map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and return
// an http.HandlerFunc (which also implements http.Handler)
// maps any paths to their corresponding URL.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYaml(yml)
	pathsToURLs := buildMap(pathURLs)
	return MapHandler(pathsToURLs, fallback), err
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathsToURLs := make(map[string]string)
	for _, pathURL := range pathURLs {
		pathsToURLs[pathURL.Path] = pathURL.URL
	}
	return pathsToURLs
}

func parseYaml(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	return pathURLs, err
}

type pathURL struct {
	Path string `yaml:"path,omitempty"`
	URL  string `yaml:"url,omitempty"`
}
