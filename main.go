package traefik_graphql_limits

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
)

const errorBodyReadResponse = `{
  "errors": [
    {
      "code": 400,
      "message": "Failed to read request body."
    }
  ]
}`

// const errorGraphqlParsingResponse = `{
//   "errors": [
//     {
//       "code": 400,
//       "message": "Failed to parse query"
//     }
//   ]
// }`

type Config struct {
	GraphQLPath string
	DepthLimit  int
	BatchLimit  int
	NodeLimit   int
}

func CreateConfig() *Config {
	return &Config{
		GraphQLPath: "/graphql",
		DepthLimit:  0,
		BatchLimit:  0,
		NodeLimit:   0,
	}
}

type GraphqlLimit struct {
	next        http.Handler
	name        string
	graphQLPath string
	depthLimit  int
	batchLimit  int
	nodeLimit   int
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &GraphqlLimit{
		next:        next,
		name:        name,
		graphQLPath: config.GraphQLPath,
		depthLimit:  config.DepthLimit,
		batchLimit:  config.BatchLimit,
		nodeLimit:   config.NodeLimit,
	}, nil
}

func respondWithJson(rw http.ResponseWriter, statusCode int, json string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(json))
	if err != nil {
		log.Printf("Error with response: %v", err)
	}
}

func (d *GraphqlLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		respondWithJson(rw, http.StatusBadRequest, errorBodyReadResponse)
		return
	}

	if req.Method == "POST" && req.URL.Path == d.graphQLPath {
		log.Printf("Checking graphql query")

		if d.depthLimit > 0 {
			log.Printf("Depth limit is set to %d", d.depthLimit)
			respondWithJson(rw, http.StatusBadRequest, errorBodyReadResponse)
		}

		if d.batchLimit > 0 {
			log.Printf("Batch limit is set to %d", d.depthLimit)
			respondWithJson(rw, http.StatusBadRequest, errorBodyReadResponse)
		}

		if d.nodeLimit > 0 {
			log.Printf("Node limit is set to %d", d.depthLimit)
			respondWithJson(rw, http.StatusBadRequest, errorBodyReadResponse)
		}
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))
	d.next.ServeHTTP(rw, req)
}
