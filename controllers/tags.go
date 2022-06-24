package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type storeRequest struct {
	Name string `json:"name" binding:"required"`
	SEO  string `json:"seo" binding:"required"`
}

func IndexTags(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}

func StoreTags(ctx *gin.Context) {
	var request storeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{})
}

func UpdateTags(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}

func DeleteTags(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{})
}
