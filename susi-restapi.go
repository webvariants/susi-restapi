package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/webvariants/susigo"
)

var cert = flag.String("cert", "cert.pem", "certificate to use")
var key = flag.String("key", "key.pem", "key to use")
var susiaddr = flag.String("susiaddr", "localhost:4000", "susiaddr to use")
var webaddr = flag.String("webaddr", ":8080", "webaddr to use")
var useHTTPS = flag.Bool("https", false, "whether to use https or not")
var mappingFile = flag.String("mapping", "endpoints.json", "the endpoint-event mapping file")

var susi *susigo.Susi

func main() {
	flag.Parse()
	s, err := susigo.NewSusi(*susiaddr, *cert, *key)
	if err != nil {
		log.Printf("Error while creating susi connection: %v", err)
		return
	}
	susi = s
	log.Println("successfully create susi connection")

	file, err := os.Open(*mappingFile)
	if err != nil {
		log.Printf("Error while opening the mapping file: %v", err)
		return
	}
	dec := json.NewDecoder(file)
	var mapping map[string]map[string]string
	err = dec.Decode(&mapping)
	if err != nil {
		log.Printf("Error while decoding the mapping file: %v", err)
		return
	}

	r := mux.NewRouter()

	for endpoint, specifier := range mapping {
		r.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			topic := specifier[r.Method]
			if _, ok := specifier[r.Method]; !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			event := susigo.Event{
				Topic: topic,
			}
			payload := make(map[string]interface{})
			for key, value := range mux.Vars(r) {
				payload[key] = value
			}
			event.Payload = payload
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				bodyPayload := make(map[string]interface{})
				payloadDecoder := json.NewDecoder(r.Body)
				err = payloadDecoder.Decode(&bodyPayload)
				if err != nil {
					w.WriteHeader(http.StatusNotAcceptable)
					fmt.Fprintf(w, "Error while parsing body: %v", err)
					return
				}
				for key, value := range bodyPayload {
					payload[key] = value
				}
				event.Payload = payload
			}
			result := make(chan interface{})
			err = susi.Publish(event, func(event *susigo.Event) {
				result <- event.Payload
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error while publishing event: %v", err)
				return
			}
			resultPayload := <-result
			enc := json.NewEncoder(w)
			enc.Encode(resultPayload)
		})
	}

	http.Handle("/", r)

	log.Printf("starting REST server on %v...", *webaddr)
	if *useHTTPS {
		log.Fatal(http.ListenAndServeTLS(*webaddr, *cert, *key, context.ClearHandler(http.DefaultServeMux)))
	} else {
		log.Fatal(http.ListenAndServe(*webaddr, context.ClearHandler(http.DefaultServeMux)))
	}
}
