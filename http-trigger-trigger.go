package main

import (
	"net/http"
	"os"
	"fmt"
	"log"
	"github.com/nu7hatch/gouuid"
	"github.com/vaughan0/go-ini"
)

var TriggerNotFoundMsg = "Trigger could not be found."

// A handler for a specific trigger.
type Trigger struct {
	// the path for which this trigger is triggered.
	path string
	// the URL that is to be called when this trigger i triggered
	outputurl string
}

// create a new trigger trigger for a specific path
func LoadHandler(outputurl, path string) Trigger {
	return Trigger{
		path: path,
		outputurl: outputurl,
	}
}

// log and trigger the specific trigger
//
// Does proper checks to make sure that the right method is from
// downstream.
func (tt Trigger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		// only supporting GET at the moment
		log.Println("[server]", r.Method, r.URL.Path, "404 (only GET allowed)")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, TriggerNotFoundMsg)
		return
	}
	if r.URL.Path != tt.path {
		// if this isn't here we'll trigger for all path
		// prefixes.
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, TriggerNotFoundMsg)
		return
	}

	loguuid, err := uuid.NewV4()
	// used to correlate server and client log lines
	var id string
	if err != nil {
		id = "???"
	} else {
		id = loguuid.String()
	}

	log.Println("[server]", id, "GET", r.URL.Path, "200")
	fmt.Fprintf(w, "Triggered.")
	go func() {
		resp, err := http.Get(tt.outputurl)
		var result string
		if err != nil {
			result = "Err: " + err.Error()
		} else {
			result = resp.Status
		}
		log.Println("[client]", id, "GET", tt.outputurl, result)
	}()
}

// Load all handlers from the configuration file.
func LoadHandlers(file ini.File) {
	root_found := false

	for path, _ := range file {
		if path == "" {
			// ignoring default section
			continue
		}

		if path == "/" {
			root_found = true
		}

		outputurl, ok := file.Get(path, "url")
		if !ok {
			panic(fmt.Sprintf("'url' missing for: %s", path))
		}

		http.Handle(path, LoadHandler(outputurl, path))
	}

	if !root_found {
		// custom 404
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("[server]", r.Method, r.URL.Path, "404")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, TriggerNotFoundMsg)
		})
	}
}

const CONFIG_FILE_ENVIRON = "TRIGGER_TRIGGER_CONFIG"

func main() {
	inifilename := os.Getenv(CONFIG_FILE_ENVIRON)
	if inifilename == "" {
		fmt.Println("Please point", CONFIG_FILE_ENVIRON, "to config file.")
		os.Exit(1)
	}

	file, err := ini.LoadFile(inifilename)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file: %s", inifilename))
	}

	LoadHandlers(file)

	listen, ok := file.Get("", "listen")
	if !ok {
		log.Print("[server] 'listen' not defined. Using fallback ':8080'")
		listen = ":8080"
	}

	log.Fatal(http.ListenAndServe(listen, nil))
}
