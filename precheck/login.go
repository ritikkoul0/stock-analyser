package precheck

import (
	"context"
	"stock-analyser/database"
	"stock-analyser/inputstructures"
	"stock-analyser/logger"
	"stock-analyser/rediscache"
)

func Login(input inputstructures.LoginInput) (bool, string, error) {
	ctx := context.Background()

	// Check if email exists in Redis
	email, err := rediscache.Client.Exists(ctx, "email:"+input.Email).Result()
	if err != nil {
		logger.Errorf("Redis error during login precheck: %v", err)
		return false, "", err
	}

	if email > 0 {
		// logger.Infof("Redis hit: Email %s found in cache", input.Email)
		return true, "", nil
	}

	// Not found in Redis, fallback to DB check
	user, err := database.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return false, "Email not registered", nil
	}

	// User found in DB, add to Redis for future logins
	rediscache.AddDataToCache(user.Username, user.Email)
	logger.Infof("Email added to redis from DB during login")

	return true, "", nil
}
