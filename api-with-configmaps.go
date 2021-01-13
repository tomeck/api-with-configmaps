package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
)

// Document - Structure to hold a document read from config and/or Provider API
type Document struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Format string `json:"format"`
	URI    string `json:"URI"`
	Tags   string `json:"tags"`
}

// downloadDoc - download the specified document's contents, and return as byte array
func downloadDoc(doc Document) ([]byte, error) {

	// Get the http object
	resp, err := http.Get(doc.URI)
	if err != nil {
		return nil, err
	}

	// Fetch the contents of its body
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// fectchDocHandler Fetch a the content of a specified document from this tenant
func fetchDocHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside fetchDocHandler")

	vars := mux.Vars(r)
	docID := vars["docId"]

	log.Printf("Requested docId %s\n", docID)

	doc := docByID(docID)

	// If we find the doc w/specified ID, marshal to JSON and return
	if doc != nil {

		contents, err := downloadDoc(*doc)

		if err == nil {

			// JTE TESTING - render any markdown doc's to HTML
			if strings.ToLower(doc.Format) == "markdown" {
				html := markdown.ToHTML(contents, nil, nil)

				w.Write([]byte(html))
			} else {
				w.Write([]byte(contents))
			}
			// END TESTING
		} else {
			log.Println("fetchDocHandler ERROR, ", err)

			// TODO - properly handle the error
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// docHandler Get a single document object from this tenant
func docHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside docHandler")

	vars := mux.Vars(r)
	docID := vars["docId"]

	log.Printf("Requested docId %s\n", docID)

	doc := docByID(docID)

	// If we find the doc w/specified ID, marshal to JSON and return
	if doc != nil {
		// Marshal the document to JSON
		js, err := json.Marshal(doc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the marshaled (JSON-ified) list of docs
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(js))
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// docsHandler Get a list of all available documents from this tenant
func docsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside docsHandler")

	// Marshal the global in-memory database of documebts
	js, err := json.Marshal(gDocuments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the marshaled (JSON-ified) list of docs
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}

// loadDocs - Load the list of documents supported by this tenant
// from ENV, which was populated from ConfigMap
func loadDocs() {

	// Fetch array of documents from ENV (which came from ConfigMap)
	docsJSON := os.Getenv("docs_json")

	json.Unmarshal([]byte(docsJSON), &gDocuments)
}

// docByID - Search the in-memory document database for specified doc
func docByID(docID string) *Document {

	var returnDoc *Document

	// Poor-man's search
	for i := range gDocuments {
		if gDocuments[i].ID == docID {
			returnDoc = &gDocuments[i]
			break
		}
	}

	return returnDoc
}

// This is our in-memory database of documents that this tenant provides
var gDocuments []Document

func main() {

	// JTE TODO remove - this is for testing webhook to kickoff build
	log.Println("Starting....[v0.0.8]")

	//
	// Load our Document database
	//
	loadDocs()

	//
	// Create our router
	//
	router := mux.NewRouter()

	//
	// Setup our resource handlers
	//
	router.HandleFunc("/docs", docsHandler).Methods("GET")
	router.HandleFunc("/docs/{docId}", docHandler).Methods("GET")
	router.HandleFunc("/docs/{docId}/content", fetchDocHandler).Methods("GET")

	//
	// Configure our server
	//
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2000)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
