package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jezzaho/goro-web/internal"
)

var progressChan = make(chan int)

const postURL = "https://api.lufthansa.com/v1/oauth/token"

func (app *Application) MockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// From parsing
	err := r.ParseForm()
	if err != nil {
		log.Println("Error during Form Parsing: ", err)
	}
	log.Println(r.Form)
	// Access the query parameters
	// Carrier Format in two letters
	carrier := r.FormValue("carrier")
	dateFrom := r.FormValue("date-from")
	dateTo := r.FormValue("date-to")
	separate := r.FormValue("separate")

	// Carrier check for Query
	var carrierNumber int

	switch carrier {
	case "LH":
		carrierNumber = 0
	case "OS":
		carrierNumber = 1
	case "LX":
		carrierNumber = 2
	case "SN":
		carrierNumber = 3
	case "EN":
		carrierNumber = 4
	default:
		carrierNumber = 0
	}

	log.Println(carrierNumber)
	dateFromSSIM := internal.DateToSSIM(dateFrom)
	dateToSSIM := internal.DateToSSIM(dateTo)
	log.Println(dateFromSSIM)
	log.Println(dateToSSIM)
	separateBool := false
	if separate == "on" {
		separateBool = true
	}
	log.Println(separateBool)

	auth, err := internal.PostForAuth(http.DefaultClient, postURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := internal.GetQueryListForAirline(carrierNumber, dateFromSSIM, dateToSSIM)
	progressChan <- 33
	data := internal.GetApiData(query, auth)
	progressChan <- 66
	data = internal.FlattenJSON(data)
	progressChan <- 100
	w.Header().Set("Content-Type", "text/csv")

	// Headerss

	currentDate := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s.csv", currentDate, carrier)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Printf("DATA: \n%v+", string(data))
	// Use the modified CreateCSVFromResponse function
	if err := internal.CreateCSVFromResponse(w, data, separateBool); err != nil {
		log.Printf("Error creating CSV: %v", err)
		http.Error(w, "Error creating CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		http.ServeFile(w, r, "static/index.html")
		return
	}
	app.fs.ServeHTTP(w, r)
}

// Live progress handler

func (app *Application) ProgressStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	clientChan := make(chan int)
	go func() {
		for progress := range progressChan {
			clientChan <- progress
			if progress >= 100 {
				close(clientChan)
				return
			}
		}
	}()
	for progress := range clientChan {
		fmt.Fprintf(w, "data: %d\n\n", progress)
		flusher.Flush()
	}

}
