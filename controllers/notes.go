package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"

	"iex/notesdot/models"
)

func IndexNotes(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}

func StoreNotes(ctx *gin.Context) {
	var request models.Note
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

	req := esapi.IndexRequest{
		Index: "notes",
		Body:  strings.NewReader(string(body)),
	}

	resp, err := req.Do(context.Background(), models.ES)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(resp.StatusCode, gin.H{"data": r})
}

func UpdateNotes(ctx *gin.Context) {
	var request models.Note
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	note := models.Note{}
	if request.Text != "" {
		note.Text = request.Text
		note.HTML = string(markdown.ToHTML([]byte(request.Text), nil, nil))
	}

	body, err := json.Marshal(note)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req := esapi.IndexRequest{
		Index:      "notes",
		DocumentID: ctx.Param("id"),
		Body:       strings.NewReader(string(body)),
	}

	resp, err := req.Do(context.Background(), models.ES)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(resp.StatusCode, gin.H{"data": r})
}

func DeleteNotes(ctx *gin.Context) {
	req := esapi.DeleteRequest{
		Index:      "notes",
		DocumentID: ctx.Param("id"),
	}

	resp, err := req.Do(context.Background(), models.ES)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Print the response status, number of results, and request duration.

	ctx.JSON(resp.StatusCode, gin.H{
		"data": r,
	})
}
