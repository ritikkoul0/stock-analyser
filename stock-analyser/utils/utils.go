package utils

type AppConfig struct {
	ServerHost string
	ServerPort string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

var Config *AppConfig

func UpdateVariables() {
	Config = &AppConfig{
		ServerHost: "localhost",
		ServerPort: "8080",
		DBHost:     "ep-ancient-bonus-a4wkt899-pooler.us-east-1.aws.neon.tech",
		DBPort:     5432,
		DBUser:     "neondb_owner",
		DBPassword: "npg_IzdP9VJ3QMXi",
		DBName:     "neondb",
		DBSSLMode:  "require",
	}
}
