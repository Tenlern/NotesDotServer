package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"iex/notesdot/controllers"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithPretty(),
		)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}

			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing the response body: %w", err)})
			} else {
				// Print the response status and error information.
				ctx.JSON(res.StatusCode, gin.H{
					"type":   e["error"].(map[string]interface{})["type"],
					"reason": e["error"].(map[string]interface{})["reason"],
				})
				return
			}
		}

		tagsRoutes := r.Group("/tags")
		{
			tagsRoutes.GET("/", controllers.IndexTags)
			tagsRoutes.GET("/", controllers.StoreTags)
			tagsRoutes.GET("/", controllers.UpdateTags)
			tagsRoutes.GET("/", controllers.DeleteTags)
		}

		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		ctx.JSON(res.StatusCode, gin.H{"data": r})
	})

	r.Run()
}
