package main

import (
	"net/http"
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
	"github.com/nu7hatch/gouuid"
	"github.com/vaughan0/go-ini"
)

var TriggerNotFoundMsg = "Trigger could not be found."

type LogUUID string

// A handler for a specific trigger.
type Trigger struct {
	// the path for which this trigger is triggered.
	path string
	// the backend URL
	outputurl string
	// channel for HTTP requests
	requestchan chan LogUUID
	// time to wait between handling each trigger
	hitDelay time.Duration
	// the maximum number of pending triggers to be handled. When
	// queue is full we start to discard incoming triggers.
	queueSize int
}

// Process HTTP requests for a trigger.
func handleHttpTriggerTriggers(t Trigger) {
	if t.outputurl == "" {
		return
	}
	for id := range t.requestchan {
		resp, err := http.Get(t.outputurl)
		var result string
		if err != nil {
			result = "Err: " + err.Error()
		} else {
			result = resp.Status
		}
		log.Println("[client]", id, "GET", t.outputurl, result)
		time.Sleep(t.hitDelay)
	}
}


// create a new trigger trigger for a specific path
func LoadHandler(path string, section ini.Section) Trigger {
	outputurl, _ := section["url"]

	sHitDelay, useHitDelay := section["hit_delay"]
	if !useHitDelay {
		sHitDelay = "0"
	}
	hitDelay, err := time.ParseDuration(sHitDelay)
	if err != nil {
		panic(fmt.Sprintf("'hit_delay', unparsable duration: %s", sHitDelay))
	}
	if hitDelay < 0 {
		panic("'hit_delay' must be positive")
	}

	sQueueSize, hasQueueSize := section["queue_size"]
	if !hasQueueSize {
		sQueueSize = "100"
	}
	queueSize, err := strconv.Atoi(sQueueSize)
	if err != nil {
		panic(fmt.Sprintf("'queue_size' not a number: %s", sQueueSize))
	}
	if queueSize < 0 {
		panic("'queue_size' must be positive")
	}

	t := Trigger{
		path: path,
		requestchan: make(chan LogUUID, queueSize),
		queueSize: queueSize,
		hitDelay: hitDelay,
		outputurl: outputurl,
	}
	go handleHttpTriggerTriggers(t)
	return t
}

const StatusTooManyRequests = 429

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

	select {
	case tt.requestchan <- LogUUID(id):
		// No rate limitting in place
		fmt.Fprintf(w, "Triggered.")
		log.Println("[server]", id, "GET", r.URL.Path, "200")
	default:
		// We were rate limitted
		w.WriteHeader(StatusTooManyRequests)
		log.Println("[server]", id, "GET", r.URL.Path, strconv.Itoa(StatusTooManyRequests))
		fmt.Fprintf(w, "Too many requests. Calm down, please.")
	}
}

// Load all handlers from the configuration file.
func LoadHandlers(file ini.File) {
	root_found := false

	for path, settings := range file {
		if path == "" {
			// ignoring default section
			continue
		}

		if path == "/" {
			root_found = true
		}

		http.Handle(path, LoadHandler(path, settings))
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
