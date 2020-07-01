package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/post", func(c *gin.Context) {

		// id := c.Query("id")
		// page := c.DefaultQuery("page", "0")
		// name := c.PostForm("name")
		// message := c.PostForm("message")
		//body, _ := ioutil.ReadAll(c.Request.Body)
		//value, dataType, offset, err := jsonparser.Get(body, "mid")
		//fmt.Printf("mid: %s; type: %T; offset: %d; err: %s", value, dataType, offset, err)

		result := getCurrentTaxiLocations()
		c.JSON(http.StatusOK, result)
	})

	router.Run(":8080")
}

func getCurrentTaxiLocations() Carparks {

	url := viperEnvVariable("Carpark")
	key := viperEnvVariable("AccountKey")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("AccountKey", key)

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(body))

	var carparks Carparks
	json.Unmarshal(body, &carparks)

	//fmt.Println(len(carparks.Carparks))

	return carparks

	// // Decoding arbitrary data: https://blog.golang.org/json
	// var f interface{}
	// err = json.Unmarshal(body, &f)

	// // iterate through the map with a range statement
	// m := f.(map[string]interface{})
	// for k, v := range m {
	// 	switch vv := v.(type) {
	// 	case string:
	// 		fmt.Println(k, "is string", vv)
	// 	case float64:
	// 		fmt.Println(k, "is float64", vv)
	// 	case []interface{}:
	// 		fmt.Println(k, "is an array:")
	// 		for i, u := range vv {
	// 			fmt.Println(i, u)
	// 		}
	// 	default:
	// 		fmt.Println(k, "is of a type I don't know how to handle")
	// 	}
	// }
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
