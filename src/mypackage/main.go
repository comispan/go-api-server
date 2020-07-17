package main

import (
	"fmt"
	"log"
	"mypackage/datasource/datasource"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/carparkLots", func(c *gin.Context) {
		url := viperEnvVariable("lta.Carpark")
		key := viperEnvVariable("lta.AccountKey")
		result := datasource.GetCarparkLots(url, key)
		c.JSON(http.StatusOK, result)
	})

	router.POST("/taxiAvailability", func(c *gin.Context) {
		url := viperEnvVariable("lta.TaxiAvailability")
		key := viperEnvVariable("lta.AccountKey")
		result := datasource.GetTaxiAvailability(url, key)
		c.JSON(http.StatusOK, result)
	})

	router.POST("/busArrival", func(c *gin.Context) {
		url := viperEnvVariable("lta.BusArrival")
		key := viperEnvVariable("lta.AccountKey")
		BusStopCode := c.Query("BusStopCode")
		ServiceNo := c.Query("ServiceNo")
		result := datasource.GetBusArrival(BusStopCode, ServiceNo, url, key)
		c.JSON(http.StatusOK, result)
	})

	t := time.Now()
	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println(formattedTime)

	router.Run(":8080")
}

// return the value of the key
func viperEnvVariable(key string) string {

	viper.SetConfigName("config")
	viper.AddConfigPath("../../")
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
