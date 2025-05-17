package routers

import (
	"stock-analyzer/stock-analyser/handlers"

	"github.com/gin-gonic/gin"
)

// SSO Login
const ssoSignupUrl = "/signup"
const ssoLoginUrl = "/login"

func SetupRouter() *gin.Engine {

	// Initialize a router with basic middleware
	baseRouter := gin.Default()
	baseMiddleware := []gin.HandlerFunc{}
	baseRouter.Use(baseMiddleware...)

	baseRouter.POST(ssoSignupUrl, handlers.UserSignup)

	// Apply middleware to router
	backendMiddleware := []gin.HandlerFunc{
		// ValidateBearerToken,
	}

	baseRouter.Use(backendMiddleware...)

	return baseRouter
}
