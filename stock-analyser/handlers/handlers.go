package handlers

import (
	"net/http"
	"stock-analyzer/stock-analyser/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(ctx *gin.Context) {
	var input SignupInput

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// Save user to database (replace this with real logic)
	err = database.SaveUser(input.Username, input.Email, string(hashedPassword))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
