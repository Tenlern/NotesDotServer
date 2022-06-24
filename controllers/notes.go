package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"

	"iex/notesdot/models"
)

type storeNoteRequest struct {
	Text string `json:"text" binding:"required"`
}

func IndexNotes(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}

func StoreNotes(ctx *gin.Context) {
	var request storeNoteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	note := models.Note{}
	note.Text = request.Text
	note.HTML = string(markdown.ToHTML([]byte(request.Text), nil, nil))

	body, err := json.Marshal(note)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req := esapi.IndexRequest{
		Index: "notes",
		Body:  strings.NewReader(string(body)),
	}
	resp, err := req.Do(context.Background(), es)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// Print the response status, number of results, and request duration.
	ctx.JSON(http.StatusAccepted, gin.H{"data": r})

	// ctx.JSON(http.StatusAccepted, gin.H{"data": req})
}

func UpdateNotes(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}

func DeleteNotes(ctx *gin.Context) {
	ctx.JSON()

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req := esapi.DeleteRequest{
		Index:      "notes",
		DocumentID: ctx.Query("id"),
	}
	resp, err := req.Do(context.Background(), es)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Print the response status, number of results, and request duration.
	ctx.JSON(http.StatusOK, gin.H{
		"data": r,
	})
}
