package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// docHandler Fetch a single document from this tenant
func docHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside docHandler")

	vars := mux.Vars(r)
	docID := vars["docId"]

	log.Printf("Requested docId %s\n", docID)

}

// docsHandler Get a list of all available documents from this tenant
func docsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside docsHandler")
}

// docsHandler Fetch the URI of the OpenAPISpec for this tenant
func oasHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside oasHandler")

	log.Println("some_var", os.Getenv("some_var"))
	log.Println("example_prop_2", os.Getenv("example_prop_2"))
	log.Println("foo", os.Getenv("foo"))
	log.Println("bar", os.Getenv("bar"))

	var jsonmap map[string]string

	x := os.Getenv("DOCS_JSON")

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(x), &jsonmap)
	log.Println("Doc1", jsonmap["doc1"])
	log.Println("Doc2", jsonmap["doc2"])
}

func main() {

	//
	// Create our router
	//
	router := mux.NewRouter()

	//
	// Setup our resource handlers
	//
	router.HandleFunc("/docs", docsHandler).Methods("GET")
	router.HandleFunc("/docs/{docId}", docHandler).Methods("GET")
	router.HandleFunc("/OAS", oasHandler).Methods("GET")

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
