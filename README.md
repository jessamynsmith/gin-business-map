# gin-business-map

Gin example project that wraps Yelp API

### Development

##### Prerequisites

1. A recent go release installed
2. A Yelp developer account and app

Fork the project on github and git clone your fork, e.g.:

    git clone https://github.com/<username>/gin-business-map.git

Ensure you have go installed on your system.

Copy .env.sample into .env and set values based on your Yelp app.

Install dependencies:

    go get .

Run server:

    go run .

Retrieve data from the API with curl (or in a browser). You can filter by different locations and terms. All query string params will be passed on to the [Yelp business search API](https://docs.developer.yelp.com/reference/v3_business_search).

curl -vk -X GET -H "Content-Type: application/json" "http://localhost:8080/api/v1/businesses/search/?location=Calgary&term=sushi&sort_by=best_match&limit=1"

curl -vk -X GET -H "Content-Type: application/json" "http://localhost:8080/api/v1/businesses/y9F-Aso24hNzbUvZNiv1MQ/"


### Deployment

This project is set up to deploy on Google cloud. You can follow the [Quickstart guide](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service), substituting this project for their example.

To deploy:

    gcloud run deploy --source . --update-env-vars YELP_API_BASE_URL=<YELP BASE URL>,YELP_API_KEY=<YOUR YELP API KEY>
