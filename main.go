package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
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

func yelpRequest(url string, queryParams url.Values, config map[string]string, jsonData any) string {
	errorMessage := ""

	authorization := fmt.Sprintf("Bearer %s", config["YELP_API_KEY"])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		errorMessage = "Unable to create request"
	} else {
		req.Header.Add("Authorization", authorization)
		req.URL.RawQuery = queryParams.Encode()

		res, doErr := http.DefaultClient.Do(req)
		if doErr != nil {
			fmt.Print(doErr.Error())
			errorMessage = "Unable to retrieve data"
		} else {
			defer res.Body.Close()

			jsonErr := json.NewDecoder(res.Body).Decode(&jsonData)
			if jsonErr != nil {
				fmt.Print(jsonErr.Error())
				errorMessage = "Unable to parse data"
			}
		}

	}
	return errorMessage
}

func RequestHandler(config map[string]string, handler func(c *gin.Context, config map[string]string)) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		handler(c, config)
	}
	return fn
}

func returnError(c *gin.Context, errorMessage string) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
}

func searchBusinesses(c *gin.Context, config map[string]string) {
	apiUrl := fmt.Sprintf("%sbusinesses/search", config["YELP_API_BASE_URL"])
	jsonData := businesses{}

	errorMessage := yelpRequest(apiUrl, c.Request.URL.Query(), config, &jsonData)

	if errorMessage != "" {
		fmt.Print(errorMessage)
		returnError(c, errorMessage)
		return
	}

	c.IndentedJSON(http.StatusOK, jsonData)
}

func businessDetails(c *gin.Context, config map[string]string) {
	businessId := c.Param("businessId")
	apiUrl := fmt.Sprintf("%sbusinesses/%s", config["YELP_API_BASE_URL"], businessId)
	jsonData := business{}

	errorMessage := yelpRequest(apiUrl, c.Request.URL.Query(), config, &jsonData)

	if errorMessage != "" {
		fmt.Print(errorMessage)
		returnError(c, errorMessage)
		return
	}

	c.IndentedJSON(http.StatusOK, jsonData)
}

func businessMap(c *gin.Context, config map[string]string) {
	c.HTML(http.StatusOK, "map.tmpl", gin.H{
		"title": "Business Map",
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	config := make(map[string]string)
	envVars := []string{"YELP_API_BASE_URL", "YELP_API_KEY"}

	for _, envVar := range envVars {
		envValue, found := os.LookupEnv(envVar)
		if !found {
			errorMessage := fmt.Sprintf("Missing environment variable: %s", envVar)
			panic(errorMessage)
		}
		config[envVar] = envValue
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/api/v1/businesses/search/", RequestHandler(config, searchBusinesses))
	router.GET("/api/v1/businesses/:businessId/", RequestHandler(config, businessDetails))
	router.GET("/", RequestHandler(config, businessMap))

	router.Run("localhost:8080")
}
