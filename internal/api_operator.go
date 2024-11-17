package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}
type ApiQuery struct {
	Airline         string
	StartDate       string
	EndDate         string
	DaysOfOperation string
	TimeMode        string
	Origin          string
	Destination     string
}

// Swaping ApiQuery Fields for faster search in-out flight A to B  - swap - B to A.
func (a *ApiQuery) Swap() {
	a.Origin, a.Destination = a.Destination, a.Origin
}

func GetApiData(queryList []ApiQuery, apiAuth Auth) []byte {
	queryResult := ""

	time.Sleep(6000 * time.Millisecond)
	for _, query := range queryList {
		time.Sleep(6000 * time.Millisecond)
		queryResult += getApiResponse(apiAuth, query)
		if queryResult == "" {
			log.Println("Empty query response before  query reverse")
		}
		// Swap query fields Origin and Destination for full result
		query.Swap()
		// Has to sleep - otherwise QPS is exceeded for Api Call
		time.Sleep(6000 * time.Millisecond)
		queryResultP2 := getApiResponse(apiAuth, query)
		if queryResultP2 == "" {
			log.Println("Empty query response after query reverse")
		}
		queryResult += queryResultP2
	}

	return []byte(queryResult)

}
func getApiResponse(auth Auth, query ApiQuery) string {

	client := http.Client{}
	getUrl := "https://api.lufthansa.com/v1/flight-schedules/flightschedules/passenger"

	queryParams := url.Values{}
	queryParams.Add("airlines", query.Airline)
	queryParams.Add("startDate", query.StartDate)
	queryParams.Add("endDate", query.EndDate)
	queryParams.Add("daysOfOperation", query.DaysOfOperation)
	queryParams.Add("timeMode", query.TimeMode)
	queryParams.Add("origin", query.Origin)
	queryParams.Add("destination", query.Destination)

	fullURL := fmt.Sprintf("%s?%s", getUrl, queryParams.Encode())

	// Perform the GET request
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Println("Error during construction of GET request: ", err.Error())
	}
	request.Header.Add("Accept", "application/json")
	authStr := "Bearer " + auth.AccessToken
	request.Header.Add("Authorization", authStr)
	log.Println("Data request send.")
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error occured during GET request from LH API: ", err.Error())
		return ""
	}	

	defer response.Body.Close()
	defer log.Println("Data request closed.")

	// Read the response Getenv
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error occured during reading response body: ", err.Error())
		return ""
	}

	// Print the response body
	body = bytes.Replace(body, []byte("]["), []byte(","), -1)
	return string(body)
}

func PostForAuth() Auth {

	postString := "https://api.lufthansa.com/v1/oauth/token"

	client := http.Client{}

	form := url.Values{}
	form.Add("client_id", os.Getenv("CLIENT_ID"))
	form.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	form.Add("grant_type", os.Getenv("GRANT_TYPE"))

	req, err := http.NewRequest("POST", postString, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Error occured during POST method: ", err.Error())
		return Auth{}
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	log.Println("Auth Request send.")
	if err != nil {
		log.Println("Error occured during request: ", err.Error())
		return Auth{}
	}
	defer resp.Body.Close()
	log.Println("Auth Request closed.")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occured during reading response: ", err.Error())
		return Auth{}
	}

	// GET REQUEST BUILDER

	var auth Auth
	err = json.Unmarshal([]byte(body), &auth)
	if err != nil || auth.AccessToken == "" {
		log.Println("Error occured during parsing the data: ", err.Error())
		return Auth{}
	}
	log.Println("Successfully retrived authentication for a request")
	return auth
}





