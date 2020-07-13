package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
