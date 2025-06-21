package handlers

import (
	"log"
	"stock-analyser/database"
	"stock-analyser/kafka"
	overview "stock-analyser/rpcclient"

	"github.com/gin-gonic/gin"
)

func GetStockDetail(ctx *gin.Context) {
	bodyBytes, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	if len(bodyBytes) > 0 {
		ctx.JSON(400, gin.H{"error": "Request body must be empty"})
		return
	}

	stockarray, err := database.GetStock(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to get data from db"})
		return
	}

	for _, stock := range stockarray {
		err := kafka.SendMessage(ctx, "stock-analyser", stock)
		if err != nil {
			log.Printf("Kafka write failed for %v: %v", stock, err)
			ctx.JSON(500, gin.H{"error": "Kafka write failed"})
			return
		}
	}
	go overview.Overview()
	ctx.JSON(200, gin.H{"message": "Stock data sent to Kafka"})
}

func DeleteStockkDetail(ctx *gin.Context) {
}

func UpdateStockDetail(ctx *gin.Context) {
}
