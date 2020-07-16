package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

// Bus services
type Services struct {
	Services []Service `json:"Services"`
}

// Bus service
type Service struct {
	ServiceNo string  `json:"ServiceNo"`
	Operator  string  `json:"Operator"`
	NextBus1  NextBus `json:"NextBus"`
	NextBus2  NextBus `json:"NextBus2"`
	NextBus3  NextBus `json:"NextBus3"`
}

// Bus details
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

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/carparkLots", func(c *gin.Context) {
		result := getCarparkLots()
		c.JSON(http.StatusOK, result)
	})

	router.POST("/taxiAvailability", func(c *gin.Context) {
		result := getTaxiAvailability()
		c.JSON(http.StatusOK, result)
	})

	router.POST("/busArrival", func(c *gin.Context) {
		BusStopCode := c.Query("BusStopCode")
		ServiceNo := c.Query("ServiceNo")
		result := getBusArrival(BusStopCode, ServiceNo)
		c.JSON(http.StatusOK, result)
	})

	t := time.Now()
	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println(formattedTime)

	router.Run(":8080")
}

func getCarparkLots() Carparks {

	url := viperEnvVariable("lta.Carpark")
	key := viperEnvVariable("lta.AccountKey")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var carparks Carparks
	json.Unmarshal(body, &carparks)

	return carparks
}

func getTaxiAvailability() Locations {

	url := viperEnvVariable("lta.TaxiAvailability")
	key := viperEnvVariable("lta.AccountKey")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var locations Locations
	json.Unmarshal(body, &locations)

	return locations
}

func getBusArrival(BusStopCode string, ServiceNo string) Services {
	url := viperEnvVariable("lta.BusArrival")
	key := viperEnvVariable("lta.AccountKey")

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

	var services Services
	json.Unmarshal(body, &services)

	return services
}

// return the value of the key
func viperEnvVariable(key string) string {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
