package routers

import (
	"stock-analyser/handlers"

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
	baseRouter.POST(ssoLoginUrl, handlers.UserLogin)

	//public ipo api's
	v1 := baseRouter.Group("/api/v1")
	{
		ipo := v1.Group("/ipo")
		{
			ipo.GET("/upcoming", handlers.GetUpcomingIPOs)
			ipo.GET("/open", handlers.GetOpenIPOs)
			ipo.GET("/closed", handlers.GetClosedIPOs)
		}

		stock := v1.Group("/stock")
		{
			stock.GET("", handlers.GetStockDetail)
			stock.DELETE("/:symbol", handlers.DeleteStockDetail)
			stock.PUT("/:symbol", handlers.UpdateStockDetail)
		}
	}

	// Apply middleware to router
	backendMiddleware := []gin.HandlerFunc{
		// ValidateBearerToken,
	}

	baseRouter.Use(backendMiddleware...)

	return baseRouter
}
