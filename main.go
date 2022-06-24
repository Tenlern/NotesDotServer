package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"iex/notesdot/controllers"
	"iex/notesdot/models"
)

func main() {
	err := models.SetupES()
	if err != nil {
		log.Fatalf("Error connecting to cluster: %s", err)
	}

	r := gin.Default()

	// r.GET("/", func(ctx *gin.Context) {
	// 	res, err := es.Search(
	// 		es.Search.WithContext(context.Background()),
	// 		es.Search.WithPretty(),
	// 	)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	// 		return
	// 	}

	// 	defer res.Body.Close()

	// 	if res.IsError() {
	// 		var e map[string]interface{}

	// 		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing the response body: %s", err)})
	// 		} else {
	// 			// Print the response status and error information.
	// 			ctx.JSON(res.StatusCode, gin.H{
	// 				"type":   e["error"].(map[string]interface{})["type"],
	// 				"reason": e["error"].(map[string]interface{})["reason"],
	// 			})
	// 			return
	// 		}
	// 	}
	// 	var r map[string]interface{}
	// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 		log.Fatalf("Error parsing the response body: %s", err)
	// 	}
	// 	// Print the response status, number of results, and request duration.
	// 	ctx.JSON(res.StatusCode, gin.H{"data": r})
	// })

	tagRoutes := r.Group("/tags")
	{
		tagRoutes.GET("/", controllers.IndexTags)
		tagRoutes.POST("/", controllers.StoreTags)
		tagRoutes.PUT("/:id", controllers.UpdateTags)
		tagRoutes.DELETE("/:id", controllers.DeleteTags)
	}
	noteRoutes := r.Group("/notes")
	{
		noteRoutes.GET("/", controllers.IndexNotes)
		noteRoutes.POST("/", controllers.StoreNotes)
		noteRoutes.PUT("/:id", controllers.UpdateNotes)
		noteRoutes.DELETE("/:id", controllers.DeleteNotes)
	}

	r.Run()
}
