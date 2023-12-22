package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"urlshort"
)

func main() {

	//get the filename to use as the source file
	filename := flag.String("filename", "", "Filename of the file to read URLs from.")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler as a default using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if strings.Contains(*filename, ".yaml") {
		// Build the YAMLHandler using the mapHandler as the
		// fallback

		//open file, read it it, parse it with YAML package
		yaml, err := os.ReadFile(*filename)
		if err != nil {
			log.Fatal(err)
		}

		yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	} else if strings.Contains(*filename, ".json") {
		json, err := os.ReadFile(*filename)

		if err != nil {
			log.Fatal(err)
		}

		jsonHandler, err := urlshort.JSONHandler(json, mapHandler)
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	} else {
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", mapHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
