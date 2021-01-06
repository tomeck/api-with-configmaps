package main

import (
	"fmt"
	"net/http"
)

// docHandler Fetch a single document from this tenant
func docHandler(w http.ResponseWriter, r *http.Request) {
}

// docsHandler Get a list of all available documents from this tenant
func docsHandler(w http.ResponseWriter, r *http.Request) {
}

// docsHandler Fetch the URI of the OpenAPISpec for this tenant
func oasHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {

	//
	// Setup our resource handlers
	//
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		docsHandler(w, r)
	})
	http.HandleFunc("/docs/{docId}", func(w http.ResponseWriter, r *http.Request) {
		docHandler(w, r)
	})
	http.HandleFunc("/OAS", func(w http.ResponseWriter, r *http.Request) {
		oasHandler(w, r)
	})

	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
