package precheck

import (
	"context"
	"fmt"
	"stock-analyser/database"
	"stock-analyser/inputstructures"
	"stock-analyser/logger"
	"stock-analyser/rediscache"
)

func Signup(input inputstructures.SignupInput) (bool, string, error) {
	ctx := context.Background()

	// Check if email exists
	email, err := rediscache.Client.Exists(ctx, "email:"+input.Email).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return false, "", err
	}
	if email > 0 {
		return false, "Email Already exists", nil
	}

	// Check if username exists
	username, err := rediscache.Client.Exists(ctx, "user:"+input.Username).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return false, "", err
	}
	if username > 0 {
		return false, "Username Already exists", nil
	}

	//This means that user is not in redis
	//check if redis has been cleared and user is in database
	if entryExists, err := database.UserExists(ctx, input.Username, input.Email); err != nil {
		return false, "", err
	} else if entryExists {
		rediscache.AddDataToCache(input.Username, input.Email)
		logger.Infof("Username and email added to Redis")
		return false, "Username or Email Already exists", nil

	}

	return true, "", nil
}
