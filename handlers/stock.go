package handlers

import (
	"log"
	"stock-analyser/database"
	"stock-analyser/kafka"
	overview "stock-analyser/rpcclient"

	"github.com/gin-gonic/gin"
)

type UpdateInput struct {
	OldSymbol string `json:"old_symbol"`
	NewSymbol string `json:"new_symbol"`
}

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

func DeleteStockDetail(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(400, gin.H{"error": "Stock symbol is required"})
		return
	}

	err := database.DeleteStock(ctx, symbol)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete stock"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Stock deleted successfully"})
}

func UpdateStockDetail(ctx *gin.Context) {
	var input UpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input format"})
		return
	}

	err := database.UpdateStock(ctx, input.OldSymbol, input.NewSymbol)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update stock"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Stock updated successfully"})
}
