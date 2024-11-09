package main

import (
	"log"
	"net/http"

	"github.com/jezzaho/goro-web/internal"
	"github.com/joho/godotenv"
)

type Application struct {
}

func main() {
	godotenv.Load("somerandomfile")

	app := Application{}
	mux := http.NewServeMux()
	mux.HandleFunc("/csv", app.MockHandler)

	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		panic(err)
	}
}
func (app *Application) MockHandler(w http.ResponseWriter, r *http.Request) {
	auth := internal.PostForAuth()
	query := internal.GetQueryListForAirline(0, "01NOV24", "05NOV24")
	data := internal.GetApiData(query, auth)
	data = internal.FlattenJSON(data)
	log.Println(string(data))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=data.csv")

	// Use the modified CreateCSVFromResponse function
	if err := internal.CreateCSVFromResponse(w, data, false, true); err != nil {
		http.Error(w, "Error creating CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
