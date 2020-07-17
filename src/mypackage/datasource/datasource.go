package datasource

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Carparks struct which contains an array of carparks
type Carparks struct {
	Carparks []Carpark `json:"value"`
}

// Carpark struct which contains the details
type Carpark struct {
	CarParkID     string `json:"CarParkID"`
	Area          string `json:"Area"`
	Development   string `json:"Development"`
	Location      string `json:"Location"`
	AvailableLots int    `json:"AvailableLots"`
	LotType       string `json:"LotType"`
	Agency        string `json:"Agency"`
}

// Locations struct which contains an array of locations (lat and long)
type Locations struct {
	Locations []Location `json:"value"`
}

// Location struct which contains the lat and long
type Location struct {
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
}

// Services struct which contains an array of services
type Services struct {
	Services []Service `json:"Services"`
}

// Service struct which contains bus service details
type Service struct {
	ServiceNo string  `json:"ServiceNo"`
	Operator  string  `json:"Operator"`
	NextBus1  NextBus `json:"NextBus"`
	NextBus2  NextBus `json:"NextBus2"`
	NextBus3  NextBus `json:"NextBus3"`
}

// NextBus struct which contains the details of the next bus
type NextBus struct {
	OriginCode       string `json:"OriginCode"`
	DestinationCode  string `json:"DestinationCode"`
	EstimatedArrival string `json:"EstimatedArrival"`
	Latitude         string `json:"Latitude"`
	Longitude        string `json:"Longitude"`
	VisitNumber      string `json:"VisitNumber"`
	Load             string `json:"Load"`
	Feature          string `json:"Feature"`
	Type             string `json:"Type"`
}

// GetCarparkLots returns the available lots
func GetCarparkLots(url string, key string) (carparks Carparks) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &carparks)

	return
}

// GetTaxiAvailability returns the available taxis locations
func GetTaxiAvailability(url string, key string) (locations Locations) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &locations)

	return
}

// GetBusArrival returns the incoming bus information
func GetBusArrival(BusStopCode string, ServiceNo string, url string, key string) (services Services) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	q := req.URL.Query()
	q.Add("BusStopCode", BusStopCode)
	if ServiceNo != "" {
		q.Add("ServiceNo", ServiceNo)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &services)

	return
}
