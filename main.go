package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type category struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type location struct {
	Address1       string      `json:"address1"`
	Address2       string      `json:"address2"`
	Address3       interface{} `json:"address3"`
	City           string      `json:"city"`
	ZipCode        string      `json:"zip_code"`
	Country        string      `json:"country"`
	State          string      `json:"state"`
	DisplayAddress []string    `json:"display_address"`
}

type business struct {
	Id           string      `json:"id"`
	Alias        string      `json:"alias"`
	Name         string      `json:"name"`
	ImageUrl     string      `json:"image_url"`
	IsClosed     bool        `json:"is_closed"`
	Url          string      `json:"url"`
	ReviewCount  int         `json:"review_count"`
	Categories   []category  `json:"categories"`
	Rating       float64     `json:"rating"`
	Coordinates  coordinates `json:"coordinates"`
	Transactions []string    `json:"transactions"`
	Price        string      `json:"price"`
	Location     location    `json:"location"`
	Phone        string      `json:"phone"`
	DisplayPhone string      `json:"display_phone"`
	Distance     float64     `json:"distance"`
}

type businesses struct {
	Businesses []business `json:"businesses"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	router := gin.Default()
	router.GET("/api/v1/businesses/search/", searchBusinesses)

	router.Run("localhost:8080")
}

func returnError(c *gin.Context, errorMessage string) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
}

func searchBusinesses(c *gin.Context) {
	apiBaseUrl, existsApiBaseUrl := os.LookupEnv("YELP_API_BASE_URL")
	apiKey, existsApiKey := os.LookupEnv("YELP_API_KEY")

	if !existsApiKey || !existsApiBaseUrl {
		returnError(c, "Server is misconfigured")
		return
	}

	url := fmt.Sprintf("%sbusinesses/search", apiBaseUrl)
	authorization := fmt.Sprintf("Bearer %s", apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		returnError(c, "Unable to create request")
		return
	}
	req.Header.Add("Authorization", authorization)
	q := c.Request.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Print(q.Encode())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		returnError(c, "Unable to retrieve data")
		return
	}
	defer res.Body.Close()

	jsonData := businesses{}

	jsonErr := json.NewDecoder(res.Body).Decode(&jsonData)
	if jsonErr != nil {
		fmt.Println(jsonErr.Error())
		returnError(c, "Unable to retrieve businesses")
	}

	c.IndentedJSON(http.StatusOK, jsonData)
}
