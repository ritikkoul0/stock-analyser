package handlers

import (
	"net/http"
	"stock-analyser/database"
	"stock-analyser/inputstructures"
	"stock-analyser/precheck"
	"stock-analyser/rediscache"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(ctx *gin.Context) {
	var input inputstructures.SignupInput

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	precheckpassed, msg, err := precheck.Signup(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error: " + err.Error()})
		return
	}

	if !precheckpassed {
		ctx.JSON(http.StatusCreated, gin.H{"message": msg})
		return
	}

	//pepperisation
	input.Password = "stock" + input.Password + "analyser"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// add username and email to redis
	rediscache.AddDataToCache(input.Username, input.Email)
	// Save user to database (replace this with real logic)
	err = database.SaveUser(ctx, input.Username, input.Email, string(hashedPassword))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func UserLogin(ctx *gin.Context) {
	var input inputstructures.LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, msg, err := precheck.Login(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	// Retrieve from DB for password
	user, err := database.GetUserByEmail(ctx, input.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Add pepper before comparing
	pepperedPassword := "stock" + input.Password + "analyser"

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pepperedPassword))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
