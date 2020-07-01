package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		getCurrentTaxiLocations()
		c.JSON(http.StatusOK, "string(resp)")
	})

	router.Run(":8080")
}

func getCurrentTaxiLocations() {

	url := "https://api.data.gov.sg/v1/transport/carpark-availability"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Decoding arbitrary data: https://blog.golang.org/json
	var f interface{}
	err = json.Unmarshal(body, &f)

	// iterate through the map with a range statement 
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
