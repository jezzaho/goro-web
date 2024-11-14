package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jezzaho/goro-web/internal"
	"github.com/joho/godotenv"
)

type Application struct {
	fs http.Handler
}

func main() {
	godotenv.Load()

	app := Application{}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	app.fs = fs

	mux.HandleFunc("/csv", app.MockHandler)
	mux.HandleFunc("/", app.IndexHandler)

	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		panic(err)
	}
}
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
	carrier := r.FormValue("carrier")
	dateFrom := r.FormValue("date-from")
	dateTo := r.FormValue("date-to")
	separate := r.FormValue("separate")

	fmt.Printf("carrier - %s, dateFrom - %s, dateTo - %s, separation - %s", carrier, dateFrom, dateTo, separate)

	auth := internal.PostForAuth()
	query := internal.GetQueryListForAirline(0, "01NOV24", "05NOV24")
	data := internal.GetApiData(query, auth)
	data = internal.FlattenJSON(data)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Use the modified CreateCSVFromResponse function
	if err := internal.CreateCSVFromResponse(w, data, false); err != nil {
		log.Printf("Error creating CSV: %v", err)
		http.Error(w, "Error creating CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}
	app.fs.ServeHTTP(w, r)
}
