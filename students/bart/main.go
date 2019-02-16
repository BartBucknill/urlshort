package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

func init() {
	// create/open boltdb and seed with path and url
	db, err := bolt.Open("url-shortener.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PathsToURLs"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		if err := b.Put([]byte("/urlshort-bolt"), []byte("https://github.com/gophercises/urlshort")); err != nil {
			return fmt.Errorf("failed to insert data: %v", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := defaultMux()

	// MapHandler with default mux as fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// YAML handler with maphandler as fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)

	// Bolt db handler with YAML handler as fallback
	db, err := bolt.Open("url-shortener.db", 0600, &bolt.Options{ReadOnly: true})
	defer db.Close()
	boltDBHandler, err := BoltDBHandler(db, yamlHandler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltDBHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
